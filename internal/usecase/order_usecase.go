package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/infra/payment"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderUsecaseContract interface {
	CreateDirectOrder(ctx context.Context, userID int64, request *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	CreateGroupBuyOrder(ctx context.Context, userID int64, request *dto.CreateGroupBuyOrderRequest) (*dto.OrderResponse, error)
	GetOrders(ctx context.Context, userID int64, page, limit int) (*dto.OrderListResponse, error)
	GetOrderByID(ctx context.Context, userID int64, orderID string) (*dto.OrderResponse, error)
	HandlePaymentNotification(ctx context.Context, notification *dto.MidtransNotification) error
	ExpireOrder(ctx context.Context, orderID string) error
}

type OrderUsecase struct {
	orderRepo        repository.OrderRepository
	paymentRepo      repository.PaymentRepository
	shippingRepo     repository.OrderShippingDetailRepository
	addressRepo      repository.AddressRepository
	variantRepo      repository.ProductVariantRepository
	stockRepo        repository.ProductVariantStockRepository
	buyerSessionRepo repository.BuyerGroupBuySessionRepository
	paymentGateway   payment.PaymentGateway
	tx               repository.TxManager
	asynqClient      *asynq.Client
	log              *logrus.Logger
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	paymentRepo repository.PaymentRepository,
	shippingRepo repository.OrderShippingDetailRepository,
	addressRepo repository.AddressRepository,
	variantRepo repository.ProductVariantRepository,
	stockRepo repository.ProductVariantStockRepository,
	buyerSessionRepo repository.BuyerGroupBuySessionRepository,
	paymentGateway payment.PaymentGateway,
	tx repository.TxManager,
	asynqClient *asynq.Client,
	log *logrus.Logger,
) OrderUsecaseContract {
	return &OrderUsecase{
		orderRepo:        orderRepo,
		paymentRepo:      paymentRepo,
		shippingRepo:     shippingRepo,
		addressRepo:      addressRepo,
		variantRepo:      variantRepo,
		stockRepo:        stockRepo,
		buyerSessionRepo: buyerSessionRepo,
		paymentGateway:   paymentGateway,
		tx:               tx,
		asynqClient:      asynqClient,
		log:              log,
	}
}

// CreateDirectOrder creates an order for direct buy flow
func (u *OrderUsecase) CreateDirectOrder(ctx context.Context, userID int64, request *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// 1. Get product variant with stock and product info
	variant, err := u.variantRepo.GetProductVariant(ctx, request.ProductVariantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Product variant not found")
		}
		return nil, err
	}

	// 2. Check stock availability
	if variant.Stock == nil || variant.Stock.CurrentStock-variant.Stock.ReservedStock < request.Quantity {
		return nil, errorx.NewBadRequestError("Insufficient stock")
	}

	// 3. Get address
	address, err := u.addressRepo.FindById(ctx, request.AddressID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Address not found")
		}
		return nil, err
	}

	// Verify address belongs to user
	if address.UserID != userID {
		return nil, errorx.NewForbiddenError("Address does not belong to user")
	}

	// 4. Calculate order amounts
	priceAtOrder := variant.Price
	subtotal := priceAtOrder * float64(request.Quantity)
	deliveryCharge := 0.0 // TODO: implement shipping calculation
	totalAmount := subtotal + deliveryCharge

	// 5. Generate order number
	orderNumber := u.generateOrderNumber()

	// 6. Create order within transaction
	var order *entity.Order
	var paymentResult *payment.VAPaymentResult

	err = u.tx.WithTransaction(ctx, func(ctx context.Context) error {
		// Create order
		order = &entity.Order{
			OrderNumber:      orderNumber,
			UserID:           userID,
			SellerID:         variant.Product.SellerID,
			ProductVariantID: variant.ID,
			Quantity:         request.Quantity,
			PriceAtOrder:     priceAtOrder,
			Subtotal:         subtotal,
			DeliveryCharge:   deliveryCharge,
			TotalAmount:      totalAmount,
			Status:           entity.OrderStatusPendingPayment,
			AddressID:        request.AddressID,
		}

		if err := u.orderRepo.Create(order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Create shipping detail (snapshot of address)
		shippingDetail := &entity.OrderShippingDetail{
			OrderID:       order.ID,
			ReceiverName:  address.ReceiverName,
			Phone:         "", // Add phone if available in address
			StreetAddress: address.StreetAddress,
			RT:            address.RT,
			RW:            address.RW,
			Village:       address.Village,
			District:      address.District,
			City:          address.City,
			Province:      address.Province,
			PostalCode:    address.PostalCode,
			Notes:         address.Notes,
		}

		if err := u.shippingRepo.Create(shippingDetail); err != nil {
			return fmt.Errorf("failed to create shipping detail: %w", err)
		}

		// Reserve stock
		if err := u.stockRepo.UpdateStock(ctx, &entity.ProductVariantStock{
			ReservedStock: variant.Stock.ReservedStock + request.Quantity,
		}, variant.ID); err != nil {
			return fmt.Errorf("failed to reserve stock: %w", err)
		}

		return nil
	})

	if err != nil {
		u.log.Errorf("Transaction failed: %v", err)
		return nil, err
	}

	// 7. Create VA payment via Midtrans (outside transaction as it's external call)
	paymentResult, err = u.paymentGateway.ChargeVA(ctx, orderNumber, int64(totalAmount), request.BankCode)
	if err != nil {
		// Rollback: update order status to cancelled and release stock
		u.orderRepo.UpdateStatus(order.ID, entity.OrderStatusCancelled)
		u.log.Errorf("Failed to create VA payment: %v", err)
		return nil, errorx.NewInternalError("Failed to create payment. Please try again.")
	}

	// 8. Save payment record
	paymentEntity := &entity.Payment{
		OrderID:              order.ID,
		Amount:               totalAmount,
		Status:               entity.PaymentStatusPending,
		PaymentMethod:        "bank_transfer",
		BankCode:             request.BankCode,
		VANumber:             paymentResult.VANumber,
		BillKey:              paymentResult.BillKey,
		BillerCode:           paymentResult.BillerCode,
		GatewayTransactionID: paymentResult.TransactionID,
		ExpiredAt:            paymentResult.ExpiredAt,
	}

	if err := u.paymentRepo.Create(paymentEntity); err != nil {
		u.log.Errorf("Failed to save payment record: %v", err)
		// Continue - the order is created, we can retry payment later
	}

	// 9. Schedule order expiration task (5 minutes)
	task, err := tasks.NewOrderExpirationTask(order.ID, orderNumber)
	if err == nil {
		_, err = u.asynqClient.Enqueue(task, asynq.ProcessIn(5*time.Minute), asynq.Queue("critical"))
		if err != nil {
			u.log.Warnf("Failed to schedule expiration task for order %s: %v", orderNumber, err)
		}
	}

	u.log.Infof("Order created: %s, VA: %s", orderNumber, paymentResult.VANumber)

	// 10. Build response
	return u.buildOrderResponse(order, paymentEntity, variant, nil), nil
}

// CreateGroupBuyOrder creates an order for group buy session
func (u *OrderUsecase) CreateGroupBuyOrder(ctx context.Context, userID int64, request *dto.CreateGroupBuyOrderRequest) (*dto.OrderResponse, error) {
	// 1. Get buyer group session
	session, err := u.buyerSessionRepo.GetSessionByID(ctx, request.BuyerGroupSessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Group buy session not found")
		}
		return nil, err
	}

	// 2. Verify user is a member of the session
	isMember := false
	for _, member := range session.Members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errorx.NewForbiddenError("You are not a member of this group buy session")
	}

	// 3. Get product variant
	variant, err := u.variantRepo.GetProductVariant(ctx, session.ProductVariantID)
	if err != nil {
		return nil, err
	}

	// 4. Get address
	address, err := u.addressRepo.FindById(ctx, request.AddressID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Address not found")
		}
		return nil, err
	}

	if address.UserID != userID {
		return nil, errorx.NewForbiddenError("Address does not belong to user")
	}

	// 5. Calculate price (may apply group buy discount later)
	priceAtOrder := variant.Price
	quantity := 1 // Group buy is typically 1 item per member
	subtotal := priceAtOrder * float64(quantity)
	deliveryCharge := 0.0
	totalAmount := subtotal + deliveryCharge

	orderNumber := u.generateOrderNumber()

	// 6. Create order
	var order *entity.Order
	var paymentResult *payment.VAPaymentResult

	err = u.tx.WithTransaction(ctx, func(ctx context.Context) error {
		order = &entity.Order{
			OrderNumber:         orderNumber,
			UserID:              userID,
			BuyerGroupSessionID: &request.BuyerGroupSessionID,
			SellerID:            variant.Product.SellerID,
			ProductVariantID:    variant.ID,
			Quantity:            quantity,
			PriceAtOrder:        priceAtOrder,
			Subtotal:            subtotal,
			DeliveryCharge:      deliveryCharge,
			TotalAmount:         totalAmount,
			Status:              entity.OrderStatusPendingPayment,
			AddressID:           request.AddressID,
		}

		if err := u.orderRepo.Create(order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Create shipping detail
		shippingDetail := &entity.OrderShippingDetail{
			OrderID:       order.ID,
			ReceiverName:  address.ReceiverName,
			StreetAddress: address.StreetAddress,
			RT:            address.RT,
			RW:            address.RW,
			Village:       address.Village,
			District:      address.District,
			City:          address.City,
			Province:      address.Province,
			PostalCode:    address.PostalCode,
			Notes:         address.Notes,
		}

		if err := u.shippingRepo.Create(shippingDetail); err != nil {
			return fmt.Errorf("failed to create shipping detail: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 7. Create VA payment
	paymentResult, err = u.paymentGateway.ChargeVA(ctx, orderNumber, int64(totalAmount), request.BankCode)
	if err != nil {
		u.orderRepo.UpdateStatus(order.ID, entity.OrderStatusCancelled)
		return nil, errorx.NewInternalError("Failed to create payment")
	}

	// 8. Save payment
	paymentEntity := &entity.Payment{
		OrderID:              order.ID,
		Amount:               totalAmount,
		Status:               entity.PaymentStatusPending,
		PaymentMethod:        "bank_transfer",
		BankCode:             request.BankCode,
		VANumber:             paymentResult.VANumber,
		BillKey:              paymentResult.BillKey,
		BillerCode:           paymentResult.BillerCode,
		GatewayTransactionID: paymentResult.TransactionID,
		ExpiredAt:            paymentResult.ExpiredAt,
	}
	u.paymentRepo.Create(paymentEntity)

	// 9. Schedule expiration
	task, _ := tasks.NewOrderExpirationTask(order.ID, orderNumber)
	u.asynqClient.Enqueue(task, asynq.ProcessIn(5*time.Minute), asynq.Queue("critical"))

	return u.buildOrderResponse(order, paymentEntity, variant, nil), nil
}

// GetOrders returns paginated orders for a user
func (u *OrderUsecase) GetOrders(ctx context.Context, userID int64, page, limit int) (*dto.OrderListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	orders, total, err := u.orderRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []dto.OrderResponse
	for _, order := range orders {
		responses = append(responses, *u.buildOrderResponse(&order, order.Payment, order.ProductVariant, order.ShippingDetail))
	}

	return &dto.OrderListResponse{
		Orders:     responses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetOrderByID returns a single order detail
func (u *OrderUsecase) GetOrderByID(ctx context.Context, userID int64, orderID string) (*dto.OrderResponse, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Order not found")
		}
		return nil, err
	}

	// Verify ownership
	if order.UserID != userID {
		return nil, errorx.NewForbiddenError("Order not found")
	}

	return u.buildOrderResponse(order, order.Payment, order.ProductVariant, order.ShippingDetail), nil
}

// HandlePaymentNotification processes Midtrans webhook
func (u *OrderUsecase) HandlePaymentNotification(ctx context.Context, notification *dto.MidtransNotification) error {
	// 1. Verify signature
	if !u.paymentGateway.VerifySignature(
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		notification.SignatureKey,
	) {
		u.log.Warnf("Invalid signature for notification: %s", notification.OrderID)
		return errorx.NewBadRequestError("Invalid signature")
	}

	// 2. Find order by order_id (which is our order_number)
	order, err := u.orderRepo.FindByOrderNumber(notification.OrderID)
	if err != nil {
		u.log.Errorf("Order not found for notification: %s", notification.OrderID)
		return errorx.NewNotFoundError("Order not found")
	}

	// 3. Get payment
	paymentEntity, err := u.paymentRepo.FindByOrderID(order.ID)
	if err != nil {
		u.log.Errorf("Payment not found for order: %s", order.ID)
		return err
	}

	// 4. Update based on transaction status
	switch notification.TransactionStatus {
	case "settlement", "capture":
		// Payment successful
		now := time.Now()
		paymentEntity.Status = entity.PaymentStatusSettlement
		paymentEntity.PaidAt = &now
		u.paymentRepo.Update(paymentEntity)

		order.Status = entity.OrderStatusPaid
		u.orderRepo.Update(order)

		u.log.Infof("Payment settled for order: %s", order.OrderNumber)

	case "pending":
		// Still waiting for payment
		u.log.Infof("Payment pending for order: %s", order.OrderNumber)

	case "expire":
		paymentEntity.Status = entity.PaymentStatusExpire
		u.paymentRepo.Update(paymentEntity)

		order.Status = entity.OrderStatusExpired
		u.orderRepo.Update(order)

		// Release reserved stock
		u.releaseStock(order.ProductVariantID, order.Quantity)

		u.log.Infof("Payment expired for order: %s", order.OrderNumber)

	case "cancel", "deny":
		paymentEntity.Status = entity.PaymentStatusCancel
		u.paymentRepo.Update(paymentEntity)

		order.Status = entity.OrderStatusCancelled
		u.orderRepo.Update(order)

		u.releaseStock(order.ProductVariantID, order.Quantity)

		u.log.Infof("Payment cancelled/denied for order: %s", order.OrderNumber)
	}

	return nil
}

// ExpireOrder is called by Asynq worker when payment timeout
func (u *OrderUsecase) ExpireOrder(ctx context.Context, orderID string) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	// Only expire if still pending
	if order.Status != entity.OrderStatusPendingPayment {
		u.log.Infof("Order %s is not pending, skipping expiration", order.OrderNumber)
		return nil
	}

	// Update order status
	order.Status = entity.OrderStatusExpired
	if err := u.orderRepo.Update(order); err != nil {
		return err
	}

	// Update payment status
	if payment, err := u.paymentRepo.FindByOrderID(order.ID); err == nil {
		payment.Status = entity.PaymentStatusExpire
		u.paymentRepo.Update(payment)
	}

	// Release stock
	u.releaseStock(order.ProductVariantID, order.Quantity)

	// Cancel transaction in Midtrans
	u.paymentGateway.CancelTransaction(ctx, order.OrderNumber)

	u.log.Infof("Order expired: %s", order.OrderNumber)
	return nil
}

// Helper functions

func (u *OrderUsecase) generateOrderNumber() string {
	// Format: ORD-YYYYMMDD-XXXXX (random suffix)
	now := time.Now()
	randomSuffix := uuid.New().String()[:8]
	return fmt.Sprintf("ORD-%s-%s", now.Format("20060102"), randomSuffix)
}

func (u *OrderUsecase) releaseStock(variantID string, quantity int) {
	stock, err := u.stockRepo.GetStockByVariantID(context.Background(), variantID)
	if err != nil {
		u.log.Warnf("Failed to get stock for release: %v", err)
		return
	}

	stock.ReservedStock -= quantity
	if stock.ReservedStock < 0 {
		stock.ReservedStock = 0
	}

	if err := u.stockRepo.UpdateStock(context.Background(), stock, variantID); err != nil {
		u.log.Warnf("Failed to release stock: %v", err)
	}
}

func (u *OrderUsecase) buildOrderResponse(order *entity.Order, payment *entity.Payment, variant *entity.ProductVariant, shipping *entity.OrderShippingDetail) *dto.OrderResponse {
	resp := &dto.OrderResponse{
		ID:             order.ID,
		OrderNumber:    order.OrderNumber,
		Status:         order.Status,
		Quantity:       order.Quantity,
		PriceAtOrder:   order.PriceAtOrder,
		Subtotal:       order.Subtotal,
		DeliveryCharge: order.DeliveryCharge,
		TotalAmount:    order.TotalAmount,
		CreatedAt:      order.CreatedAt,
	}

	// Add payment details
	if payment != nil {
		resp.Payment = &dto.PaymentDetailResponse{
			ID:         payment.ID,
			BankCode:   payment.BankCode,
			VANumber:   payment.VANumber,
			BillKey:    payment.BillKey,
			BillerCode: payment.BillerCode,
			Amount:     payment.Amount,
			Status:     payment.Status,
			ExpiredAt:  payment.ExpiredAt,
			PaidAt:     payment.PaidAt,
		}
	}

	// Add product details
	if variant != nil {
		resp.Product = &dto.OrderProductResponse{
			VariantID:   variant.ID,
			VariantName: variant.Name,
		}
		if variant.Product != nil {
			resp.Product.ProductID = variant.Product.ID
			resp.Product.ProductName = variant.Product.Title
			if len(variant.Product.ProductImages) > 0 {
				resp.Product.ImageURL = variant.Product.ProductImages[0].ImageURL
			}
		}
	}

	// Add shipping details
	if shipping != nil {
		resp.ShippingDetail = &dto.ShippingDetailResponse{
			ReceiverName:  shipping.ReceiverName,
			Phone:         shipping.Phone,
			StreetAddress: shipping.StreetAddress,
			Village:       shipping.Village,
			District:      shipping.District,
			City:          shipping.City,
			Province:      shipping.Province,
			PostalCode:    shipping.PostalCode,
		}
	}

	// Add seller details
	if order.Seller != nil {
		resp.Seller = &dto.OrderSellerResponse{
			ID:       order.Seller.ID,
			ShopName: order.Seller.StoreName,
		}
	}

	return resp
}

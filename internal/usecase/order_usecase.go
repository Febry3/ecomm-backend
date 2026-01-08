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

func (u *OrderUsecase) CreateDirectOrder(ctx context.Context, userID int64, request *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	variant, err := u.variantRepo.GetProductVariant(ctx, request.ProductVariantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Product variant not found")
		}
		return nil, err
	}

	if variant.Stock == nil || variant.Stock.CurrentStock-variant.Stock.ReservedStock < request.Quantity {
		return nil, errorx.NewBadRequestError("Insufficient stock")
	}

	address, err := u.addressRepo.FindById(ctx, request.AddressID, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.log.Error("[Order Usecase] address not found: ", err)
			return nil, errorx.NewNotFoundError("Address not found")
		}
		u.log.Error("[Order Usecase] failed to get address: ", err)
		return nil, err
	}

	if address.UserID != userID {
		u.log.Error("[Order Usecase] address does not belong to user")
		return nil, errorx.NewForbiddenError("Address does not belong to user")
	}

	priceAtOrder := variant.Price
	subtotal := priceAtOrder * float64(request.Quantity)
	deliveryCharge := 0.0 // free for now (im too lazy to implement it)
	totalAmount := subtotal + deliveryCharge

	orderNumber := u.generateOrderNumber()

	var order *entity.Order
	var paymentResult *payment.VAPaymentResult

	err = u.tx.WithTransaction(ctx, func(ctx context.Context) error {
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

	paymentResult, err = u.paymentGateway.ChargeVA(ctx, orderNumber, int64(totalAmount), request.BankCode, nil)
	if err != nil {
		u.orderRepo.UpdateStatus(order.ID, entity.OrderStatusCancelled)
		u.log.Errorf("Failed to create VA payment: %v", err)
		return nil, errorx.NewInternalError("Failed to create payment. Please try again.")
	}

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

	task, err := tasks.NewOrderExpirationTask(order.ID, orderNumber)
	if err == nil {
		_, err = u.asynqClient.Enqueue(task, asynq.ProcessIn(5*time.Minute), asynq.Queue("critical"))
		if err != nil {
			u.log.Warnf("Failed to schedule expiration task for order %s: %v", orderNumber, err)
		}
	}

	u.log.Infof("Order created: %s, VA: %s", orderNumber, paymentResult.VANumber)

	return u.buildOrderResponse(order, paymentEntity, variant, nil), nil
}

func (u *OrderUsecase) CreateGroupBuyOrder(ctx context.Context, userID int64, request *dto.CreateGroupBuyOrderRequest) (*dto.OrderResponse, error) {
	session, err := u.buyerSessionRepo.GetSessionByID(ctx, request.BuyerGroupSessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Group buy session not found")
		}
		return nil, err
	}

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

	variant, err := u.variantRepo.GetProductVariant(ctx, session.ProductVariantID)
	if err != nil {
		return nil, err
	}

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

	priceAtOrder := variant.Price
	quantity := 1 // Group buy is typically 1 item per member
	subtotal := priceAtOrder * float64(quantity)
	deliveryCharge := 0.0
	totalAmount := subtotal + deliveryCharge

	orderNumber := u.generateOrderNumber()

	var order *entity.Order
	var paymentResult *payment.VAPaymentResult

	err = u.tx.WithTransaction(ctx, func(ctx context.Context) error {
		order = &entity.Order{
			OrderNumber:         orderNumber,
			UserID:              userID,
			BuyerGroupSessionID: &session.ID,
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

	paymentResult, err = u.paymentGateway.ChargeVA(ctx, orderNumber, int64(totalAmount), request.BankCode, &session.ExpiresAt)
	if err != nil {
		u.orderRepo.UpdateStatus(order.ID, entity.OrderStatusCancelled)
		return nil, errorx.NewInternalError("Failed to create payment")
	}

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

	task, _ := tasks.NewOrderExpirationTask(order.ID, orderNumber)
	u.asynqClient.Enqueue(task, asynq.ProcessIn(5*time.Minute), asynq.Queue("critical"))

	return u.buildOrderResponse(order, paymentEntity, variant, nil), nil
}

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

func (u *OrderUsecase) GetOrderByID(ctx context.Context, userID int64, orderID string) (*dto.OrderResponse, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewNotFoundError("Order not found")
		}
		return nil, err
	}

	if order.UserID != userID {
		return nil, errorx.NewForbiddenError("Order not found")
	}

	return u.buildOrderResponse(order, order.Payment, order.ProductVariant, order.ShippingDetail), nil
}

func (u *OrderUsecase) HandlePaymentNotification(ctx context.Context, notification *dto.MidtransNotification) error {
	if !u.paymentGateway.VerifySignature(
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		notification.SignatureKey,
	) {
		u.log.Warnf("Invalid signature for notification: %s", notification.OrderID)
		return errorx.NewBadRequestError("Invalid signature")
	}

	order, err := u.orderRepo.FindByOrderNumber(notification.OrderID)
	if err != nil {
		u.log.Errorf("Order not found for notification: %s", notification.OrderID)
		return errorx.NewNotFoundError("Order not found")
	}

	paymentEntity, err := u.paymentRepo.FindByOrderID(order.ID)
	if err != nil {
		u.log.Errorf("Payment not found for order: %s", order.ID)
		return err
	}

	switch notification.TransactionStatus {
	case "settlement", "capture":
		now := time.Now()
		paymentEntity.Status = entity.PaymentStatusSettlement
		paymentEntity.PaidAt = &now
		u.paymentRepo.Update(paymentEntity)

		order.Status = entity.OrderStatusPaid
		u.orderRepo.Update(order)

		if err := u.deductStockOnPayment(ctx, order.ProductVariantID, order.Quantity); err != nil {
			u.log.Errorf("Failed to deduct stock for order %s: %v", order.OrderNumber, err)
		}

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

func (u *OrderUsecase) ExpireOrder(ctx context.Context, orderID string) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != entity.OrderStatusPendingPayment {
		u.log.Infof("Order %s is not pending, skipping expiration", order.OrderNumber)
		return nil
	}

	order.Status = entity.OrderStatusExpired
	if err := u.orderRepo.Update(order); err != nil {
		return err
	}

	if payment, err := u.paymentRepo.FindByOrderID(order.ID); err == nil {
		payment.Status = entity.PaymentStatusExpire
		u.paymentRepo.Update(payment)
	}

	u.releaseStock(order.ProductVariantID, order.Quantity)

	u.paymentGateway.CancelTransaction(ctx, order.OrderNumber)

	u.log.Infof("Order expired: %s", order.OrderNumber)
	return nil
}

func (u *OrderUsecase) generateOrderNumber() string {
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

// deductStockOnPayment decrements current_stock and reserved_stock atomically with optimistic locking.
// This is called when payment is successful to finalize the stock deduction.
func (u *OrderUsecase) deductStockOnPayment(ctx context.Context, variantID string, quantity int) error {
	// Get current stock to retrieve version
	stock, err := u.stockRepo.GetStockByVariantID(ctx, variantID)
	if err != nil {
		return fmt.Errorf("failed to get stock: %w", err)
	}

	// Retry loop for optimistic lock conflict (max 3 attempts)
	for attempt := 0; attempt < 3; attempt++ {
		err = u.stockRepo.DeductStockWithVersion(ctx, variantID, quantity, stock.Version)
		if err == nil {
			u.log.Infof("Stock deducted for variant %s: qty=%d, version=%d", variantID, quantity, stock.Version)
			return nil
		}

		// If version mismatch, retry with fresh version
		stock, err = u.stockRepo.GetStockByVariantID(ctx, variantID)
		if err != nil {
			return fmt.Errorf("failed to get stock for retry: %w", err)
		}
		u.log.Warnf("Optimistic lock conflict for variant %s, retrying (attempt %d)", variantID, attempt+1)
	}

	return fmt.Errorf("failed to deduct stock after retries: version conflict")
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

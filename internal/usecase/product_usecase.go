package usecase

import (
	"context"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type ProductUsecaseContract interface {
	CreateProduct(ctx context.Context, request dto.CreateProductRequest, sellerID int64) (*entity.Product, error)
	GetAllProductsForBuyer(ctx context.Context) ([]entity.Product, error)
	GetProductForBuyer(ctx context.Context, productID string) (*entity.Product, error)
	GetAllProductsForSeller(ctx context.Context, sellerId int64) ([]entity.Product, int, int, float64, int, error)
	GetProductForSeller(ctx context.Context, productID string, sellerId int64) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, product dto.UpdateProductRequest, productID string, sellerID int64) (*dto.ProductResponse, error)
	GetAllCategories(ctx context.Context) ([]dto.CategoryResponse, error)
}

type ProductUsecase struct {
	productRepo  repository.ProductRepository
	variantRepo  repository.ProductVariantRepository
	stockRepo    repository.ProductVariantStockRepository
	sellerRepo   repository.SellerRepository
	categoryRepo repository.CategoryRepository
	tx           repository.TxManager
	log          *logrus.Logger
}

func NewProductUsecase(
	productRepo repository.ProductRepository,
	variantRepo repository.ProductVariantRepository,
	stockRepo repository.ProductVariantStockRepository,
	sellerRepo repository.SellerRepository,
	categoryRepo repository.CategoryRepository,
	tx repository.TxManager,
	log *logrus.Logger,
) ProductUsecaseContract {
	return &ProductUsecase{
		productRepo:  productRepo,
		variantRepo:  variantRepo,
		stockRepo:    stockRepo,
		sellerRepo:   sellerRepo,
		categoryRepo: categoryRepo,
		tx:           tx,
		log:          log,
	}
}

func (p *ProductUsecase) GetAllCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := p.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		p.log.Error("[ProductUsecase] GetAllCategories Error: ", err)
		return nil, err
	}
	return dto.ToCategoryResponse(categories), nil
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, request dto.CreateProductRequest, sellerID int64) (*entity.Product, error) {
	if err := validator.New().Struct(request); err != nil {
		p.log.Errorf("[ProductUsecase] Validate Product Error: %v", err.Error())
		return nil, err
	}

	product := &entity.Product{
		SellerID:    sellerID,
		Title:       request.Title,
		Slug:        request.Slug,
		Description: datatypes.JSON(request.Description),
		Badge:       request.Badge,
		IsActive:    request.IsActive,
	}

	err := p.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := p.productRepo.CreateProduct(txCtx, product); err != nil {
			p.log.Errorf("[ProductUsecase] Create Product Error: %v", err)
			return err
		}

		for _, v := range request.Variants {
			variant := &entity.ProductVariant{
				ProductID: product.ID,
				Sku:       v.Sku,
				Name:      v.Name,
				Price:     v.Price,
				IsActive:  v.IsActive,
			}

			if err := p.variantRepo.CreateProductVariant(txCtx, variant); err != nil {
				p.log.Errorf("[ProductUsecase] Create Variant Error (SKU: %s): %v", v.Sku, err)
				return err
			}

			stock := &entity.ProductVariantStock{
				ProductVariantID:  variant.ID,
				CurrentStock:      0,
				ReservedStock:     0,
				LowStockThreshold: 5,
			}

			if v.Stock != nil {
				stock.CurrentStock = v.Stock.CurrentStock
				stock.ReservedStock = v.Stock.ReservedStock
				if v.Stock.LowStockThreshold > 0 {
					stock.LowStockThreshold = v.Stock.LowStockThreshold
				}
			}

			if err := p.stockRepo.CreateStock(txCtx, stock); err != nil {
				p.log.Errorf("[ProductUsecase] Create Stock Error (Variant: %s): %v", variant.ID, err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		p.log.Errorf("[ProductUsecase] CreateProduct transaction failed: %v", err)
		return nil, err
	}

	return product, nil
}

func (p *ProductUsecase) GetAllProductsForBuyer(ctx context.Context) ([]entity.Product, error) {
	return p.productRepo.GetProductsForBuyer(ctx)
}

func (p *ProductUsecase) GetProductForBuyer(ctx context.Context, productID string) (*entity.Product, error) {
	return p.productRepo.GetProductForBuyer(ctx, productID)
}

func (p *ProductUsecase) GetAllProductsForSeller(ctx context.Context, sellerId int64) ([]entity.Product, int, int, float64, int, error) {
	products, err := p.productRepo.GetProductsForSeller(ctx, sellerId)
	if err != nil {
		p.log.Errorf("[ProductUsecase] GetProductsForSeller Error: %v", err)
		return nil, 0, 0, 0, 0, err
	}

	countVariant := 0
	totalStock := 0
	totalInventoryValue := 0.0
	totalStockAlert := 0

	for _, product := range products {
		countVariant += len(product.Variants)
		for _, variant := range product.Variants {
			totalStock += variant.Stock.CurrentStock
			totalInventoryValue += variant.Price * float64(variant.Stock.CurrentStock)
			if variant.Stock.CurrentStock < variant.Stock.LowStockThreshold {
				totalStockAlert++
			}
		}
	}
	return products, countVariant, totalStock, totalInventoryValue, totalStockAlert, nil
}

func (p *ProductUsecase) GetProductForSeller(ctx context.Context, productID string, sellerId int64) (*dto.ProductResponse, error) {
	product, err := p.productRepo.GetProductForSeller(ctx, productID, sellerId)
	if err != nil {
		return nil, err
	}

	productVariants, err := p.variantRepo.GetProductVariants(ctx, productID)
	if err != nil {
		return nil, err
	}

	return dto.ToProductResponse(product, productVariants), nil
}

func (p *ProductUsecase) DeleteProductVariant(ctx context.Context, productVariantID string) error {
	return p.variantRepo.DeleteProductVariant(ctx, productVariantID)
}

func (p *ProductUsecase) UpdateProduct(ctx context.Context, product dto.UpdateProductRequest, productID string, sellerID int64) (*dto.ProductResponse, error) {
	if err := validator.New().Struct(product); err != nil {
		p.log.Errorf("[ProductUsecase] Validate Product Error: %v", err.Error())
		return nil, err
	}

	productEntity := &entity.Product{
		ID:          productID,
		Title:       product.Title,
		Slug:        product.Slug,
		Description: datatypes.JSON(product.Description),
		Badge:       product.Badge,
		IsActive:    *product.IsActive,
	}

	var updatedVariant []entity.ProductVariant
	err := p.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := p.productRepo.UpdateProductForSeller(txCtx, productEntity, productID, sellerID); err != nil {
			p.log.Errorf("[ProductUsecase] Update Product Error: %v", err)
			return err
		}

		for _, v := range product.ProductVariants {
			variant := &entity.ProductVariant{
				ProductID: productID,
				Sku:       v.Sku,
				Name:      v.Name,
				Price:     v.Price,
				IsActive:  v.IsActive,
			}

			if v.ID == "" {
				if err := p.variantRepo.CreateProductVariant(txCtx, variant); err != nil {
					p.log.Errorf("[ProductUsecase] Create Variant Error (SKU: %s): %v", v.Sku, err)
					return err
				}
				if err := p.stockRepo.CreateStock(txCtx, &entity.ProductVariantStock{
					ProductVariantID:  variant.ID,
					CurrentStock:      v.Stock.CurrentStock,
					ReservedStock:     v.Stock.ReservedStock,
					LowStockThreshold: v.Stock.LowStockThreshold,
				}); err != nil {
					p.log.Errorf("[ProductUsecase] Create Stock Error (Variant: %s): %v", variant.ID, err)
					return err
				}
				break
			}

			if err := p.variantRepo.UpdateProductVariant(txCtx, variant, v.ID); err != nil {
				p.log.Errorf("[ProductUsecase] Update Variant Error (SKU: %s): %v", v.Sku, err)
				return err
			}

			stock := &entity.ProductVariantStock{
				ProductVariantID:  v.ID,
				CurrentStock:      0,
				ReservedStock:     0,
				LowStockThreshold: 5,
			}

			if v.Stock != nil {
				stock.CurrentStock = v.Stock.CurrentStock
				stock.ReservedStock = v.Stock.ReservedStock
				if v.Stock.LowStockThreshold > 0 {
					stock.LowStockThreshold = v.Stock.LowStockThreshold
				}
			}

			if err := p.stockRepo.UpdateStock(txCtx, stock, v.ID); err != nil {
				p.log.Errorf("[ProductUsecase] Update Stock Error (Variant: %s): %v", v.ID, err)
				return err
			}
			variant.Stock = stock
			updatedVariant = append(updatedVariant, *variant)
		}

		return nil
	})

	if err != nil {
		p.log.Errorf("[ProductUsecase] UpdateProduct transaction failed: %v", err)
		return nil, err
	}

	return dto.ToProductResponse(productEntity, updatedVariant), nil
}

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sample-crud/internal/domain"
	customerrors "sample-crud/pkg/errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductService struct {
	productRepository domain.ProductRepository
	redisClient       *redis.Client
}

func (p ProductService) CreateProduct(ctx context.Context, name string) (uint, error) {
	log := logger.GetLogger(ctx)
	log.SInfo("Starting product creation with name : %s", name)
	id, err := p.productRepository.Create(ctx, &domain.Product{Name: name})
	if err != nil {
		log.Error("Fail to create product", zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (p ProductService) FindByID(ctx context.Context, id uint) (*domain.ProductInfo, error) {
	log := logger.GetLogger(ctx)
	log.SInfo("Starting finding product with id : %d", id)

	cacheKey := fmt.Sprintf("sample_crud:product#%d", id)
	cacheData, err := p.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		log.SInfo("Product cache found product with id : %d", id)
		var product domain.ProductInfo
		if err := json.Unmarshal([]byte(cacheData), &product); err == nil {
			return &product, nil
		}
		log.SWarn("Unmarshal product cache failed")
	} else if !errors.Is(err, redis.Nil) {
		log.SWarn("Redis error : %v", err)
	}

	log.SInfo("Product cache not found with id : %d", id)
	product, err := p.productRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.SInfo("Record not found for product with id : %d", id)
			return nil, customerrors.NewNotFoundError("Product not found", err.Error())
		}
		log.Error("Fail to find product by id", zap.Error(err))
		return nil, err
	}
	productInfo := domain.ProductInfo{
		ID:   product.ID,
		Name: product.Name,
	}
	saveProductCache(ctx, p.redisClient, product, log)
	return &productInfo, nil
}

func (p ProductService) UpdateProduct(ctx context.Context, id uint, name string) error {
	log := logger.GetLogger(ctx)
	log.SInfo("Starting product update with id : %d and name : %s", id, name)

	existingProduct, err := p.productRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.SInfo("Record not found for product with id : %d", id)
			return customerrors.NewNotFoundError("Product not found", err.Error())
		}
		log.Error("Fail to find product by id for update", zap.Error(err))
		return err
	}

	existingProduct.Name = name

	if err := p.productRepository.Update(ctx, existingProduct); err != nil {
		log.Error("Fail to update product", zap.Error(err))
		return err
	}

	invalidateProductCache(ctx, p.redisClient, id, log)

	log.SInfo("Product updated successfully with id : %d", id)
	return nil
}

func (p ProductService) DeleteProduct(ctx context.Context, id uint) error {
	log := logger.GetLogger(ctx)
	log.SInfo("Starting product deletion with id : %d", id)

	rowsAffected, err := p.productRepository.Delete(ctx, id)
	if err != nil {
		log.Error("Fail to delete product", zap.Error(err))
		return err
	}

	if rowsAffected == 0 {
		log.SInfo("Record not found for product with id : %d", id)
		return customerrors.NewNotFoundError("Product not found", "no rows affected")
	}

	invalidateProductCache(ctx, p.redisClient, id, log)

	log.SInfo("Product deleted successfully with id : %d", id)
	return nil
}

func getProductCacheKey(id uint) string {
	return fmt.Sprintf("sample_crud:product#%d", id)
}

func saveProductCache(ctx context.Context, redisClient *redis.Client, product *domain.Product, log *logger.Logger) {
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Warn("Fail to marshal product", zap.Error(err))
		return
	}
	err = redisClient.Set(ctx, getProductCacheKey(product.ID), string(productJSON), time.Minute*30).Err()
	if err != nil {
		log.Warn("Fail to save product cache", zap.Error(err))
		return
	}
	log.SInfo("Product cache saved successfully")
}

func invalidateProductCache(ctx context.Context, redisClient *redis.Client, id uint, log *logger.Logger) {
	cacheKey := getProductCacheKey(id)
	err := redisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Warn("Fail to delete product cache", zap.Error(err))
	}
}

func NewProductService(productRepository domain.ProductRepository, redisClient *redis.Client) *ProductService {
	return &ProductService{
		productRepository: productRepository,
		redisClient:       redisClient,
	}
}

package repo

import (
	"context"
	"sample-crud/internal/domain"
	"time"

	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func (g GormProductRepository) Create(ctx context.Context, product *domain.Product) (uint, error) {
	var now = time.Now()
	product.ID = 0
	product.CreatedAt = &now
	product.UpdatedAt = &now
	if err := g.db.WithContext(ctx).Create(product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (g GormProductRepository) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	if err := g.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (g GormProductRepository) Update(ctx context.Context, product *domain.Product) error {
	var now = time.Now()
	product.UpdatedAt = &now

	if err := g.db.WithContext(ctx).Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (g GormProductRepository) Delete(ctx context.Context, id uint) (int64, error) {
	result := g.db.WithContext(ctx).Delete(&domain.Product{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

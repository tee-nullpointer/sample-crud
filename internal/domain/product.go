package domain

import (
	"time"
)

type Product struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	CreatedAt *time.Time // với gorm luôn nên để con trỏ để có giá trị nil, nếu không sẽ insert giờ mặc định
	UpdatedAt *time.Time
}

func (p Product) TableName() string {
	return "sample.products" // implement interface Tabler chỉ định schema và table name
}

type ProductCreation struct {
	Name string `json:"name" binding:"required,min=3"`
}

type ProductUpdate struct {
	Name string `json:"name" binding:"required,min=3"`
}

type ProductInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// internal/models/product.go
package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Orders   []Order `gorm:"foreignKey:ProductID"` // Relation
}

type Order struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

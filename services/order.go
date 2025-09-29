package services

import (
	"erp-system/models"

	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func (s *OrderService) CreateOrder(userID uint, productID uint, qty int) (*models.Order, error) {
	var product models.Product
	if err := s.DB.First(&product, productID).Error; err != nil {
		return nil, err
	}
	if product.Quantity < qty {
		return nil, gorm.ErrRecordNotFound // Or custom error
	}

	// Transaction: Deduct stock + create order
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order models.Order
	order.UserID = userID
	order.ProductID = productID
	order.Quantity = qty
	order.Total = float64(qty) * product.Price

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	product.Quantity -= qty
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &order, tx.Commit().Error
}

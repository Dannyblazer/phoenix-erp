package controllers

import (
	"erp-system/initializers"
	"erp-system/models"
	"erp-system/services"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productInputSerializer struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Quantity int     `json:"quantity" binding:"required,gte=0"`
}

type productSerializer struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type orderInputSerializer struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func ProductCreate(c *gin.Context) {
	// Get user from jwt auth
	userID, ok := c.Get("userID")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No user associated",
		})
		return
	}
	// Search userID in DB
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Associated user not found",
		})
		return
	}

	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Only Admins can create products",
		})
		return
	}

	var req productInputSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Object",
		})
		return
	}

	product := models.Product{
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	}

	if err := initializers.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
	})
}

func ProductList(c *gin.Context) {
	var products []models.Product
	if err := initializers.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func ProductGet(c *gin.Context) {
	// Get and validate productID
	idStr := c.Param("id")
	fmt.Printf("Here's the product ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	// fetch product
	var product models.Product
	if err := initializers.DB.Select("id", "name", "quantity", "price", "created_at", "updated_at").
		First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product Not found",
			})
			return
		}
		log.Printf("Failed to fetch product %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch product",
		})
		return
	}

	response := productSerializer{
		ID:         product.ID,
		Name:       product.Name,
		Quantity:   product.Quantity,
		Price:      product.Price,
		CreatedAt:  product.CreatedAt,
		Updated_At: product.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"product": response,
	})
}

func OrderCreate(c *gin.Context) {
	// Get UserID
	fmt.Println("Order service Reached")
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No User associated",
		})
		return
	}
	// Get user account

	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	// Parse and validate request body
	orderSvc := &services.OrderService{DB: initializers.DB}

	var req orderInputSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Body",
		})
		return
	}

	// Create Order
	order, err := orderSvc.CreateOrder(uint(user.ID), uint(req.ProductID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"order": order,
	})
}

func ProductUpdate(c *gin.Context) {
	// Get and validate productID
	idStr := c.Param("id")
	fmt.Printf("Here'{{local}}orders/s the product ID: %s", idStr)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	var req productInputSerializer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Body",
		})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No user associated",
		})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Associated user not found",
		})
		return
	}

	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Only Admins can create products",
		})
		return
	}

	var product models.Product
	if err := initializers.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	// Update product fields
	product.Name = req.Name
	product.Price = req.Price
	product.Quantity = req.Quantity

	// Save updated product
	if err := initializers.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update product",
		})
		return
	}

	response := productSerializer{
		ID:         product.ID,
		Name:       product.Name,
		Quantity:   product.Quantity,
		Price:      product.Price,
		CreatedAt:  product.CreatedAt,
		Updated_At: product.UpdatedAt,
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Update Successful",
		"product": response,
	})

}

func OrderList(c *gin.Context) {
	var orders []models.Order
	if err := initializers.DB.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

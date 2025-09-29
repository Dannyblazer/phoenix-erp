package controllers

import (
	"erp-system/initializers"
	"erp-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productInput struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Quantity int     `json:"quantity" binding:"required,gte=0"`
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

	var req productInput
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

}

func OrderCreate(c *gin.Context) {

}

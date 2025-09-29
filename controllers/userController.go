package controllers

import (
	"erp-system/initializers"
	"erp-system/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountBody struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Role     string `json:"role" binding:"required,max=20"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type AccountResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserCreate(c *gin.Context) {
	var req AccountBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Object",
		})
		return
	}

	// Validate user role
	req.Role = strings.ToLower(strings.TrimSpace((req.Role)))
	if req.Role != "admin" && req.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User role specified",
		})
		return
	}

	// Normalize and hash password
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to hash password",
		})
		return
	}

	// Check and validate email uniqueness
	var existingAccount models.User
	if err := initializers.DB.Where("email = ?", req.Email).First(&existingAccount).Error; err == nil {
		log.Printf("Email Already used")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email Already Exists",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("Error Checking existing account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error checking account",
		})
		return
	}

	account := models.User{
		Email:    req.Email,
		Role:     req.Role,
		Password: string(hash),
	}

	if err := initializers.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error Creating account",
		})
		return
	}

	response := AccountResponse{
		ID:        account.ID,
		Email:     account.Email,
		Role:      account.Role,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created",
		"user":    response,
	})

}


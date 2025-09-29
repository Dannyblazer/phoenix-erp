// internal/models/user.go
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string // Hash in production
	Role     string // "admin", "user"
}

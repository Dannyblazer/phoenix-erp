package initializers

import "erp-system/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
}

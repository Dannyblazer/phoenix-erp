package main

import (
	"erp-system/controllers"
	"erp-system/initializers"
	"erp-system/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	initializers.LoadEnvVariables() // Load env variables
	initializers.ConnectDB()        // connect to DB
	initializers.SyncDB()           // Sync DB
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := gin.Default()

	router.POST("account/login/", controllers.UserLogin)
	router.POST("account/create/", controllers.UserCreate)
	router.POST("products/", middleware.RequireAuth, controllers.ProductCreate)
	router.GET("products/:id", middleware.RequireAuth, controllers.ProductGet)
	router.PUT("products/:id", middleware.RequireAuth, controllers.ProductUpdate)
	router.GET("products/list/", middleware.RequireAuth, controllers.ProductList)
	logger.Info("Starting Server on port 8000")
	router.Run()
}

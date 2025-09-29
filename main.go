package main

import (
	"erp-system/initializers"
	"net/http"

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

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ERP System Ready!",
		})
	})

	logger.Info("Starting Server on port 8000")
	r.Run()
}

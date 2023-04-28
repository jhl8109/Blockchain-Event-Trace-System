package main

import (
	"Blockchain-Event-Trace-System/db"
	"Blockchain-Event-Trace-System/gateway"
	"github.com/gin-gonic/gin"
)

func main() {
	gateway.Connect()

	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	db.ConnectionDB()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//router.GET("/connect", func(c *gin.Context) {
	//	gateway.Connect()
	//	c.JSON(200, gin.H{
	//		"message": "connected",
	//	})
	//})
	router.POST("/asset", gateway.CreateAsset)
	router.PUT("/asset", gateway.UpdateAsset)
	router.POST("/tx", gateway.TransferAsset)
	router.DELETE("/asset", gateway.DeleteAsset)
	router.GET("/asset", gateway.GetAllAssets)
	router.POST("/tt", gateway.GetAsset)
	return router
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/princesp/go-jwt/initializer"
)

func init() {
	initializer.LoadEnvInitializer()
	initializer.ConnectToDB()
	initializer.SyncDatabase()
}
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
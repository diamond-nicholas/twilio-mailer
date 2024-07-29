package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nicholas/go-jwt/controllers"
	"github.com/nicholas/go-jwt/initializers"
	"github.com/nicholas/go-jwt/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

		r.POST("/signup", controllers.Signup)
		r.POST("/login", controllers.Login)
		r.GET("/validate", middleware.RequireAuth,  controllers.Validate)
		r.POST("/upload", controllers.UploadAndSendEmails)
	r.Run()
}

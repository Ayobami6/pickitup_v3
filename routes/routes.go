package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RootRoute(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		fmt.Println("Testing debug mode another")
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the API!",
			"version":  "3.0",
            "author": "Ayobami Alaran",
            "contact": "https://github.com/Ayobami6/pickitup_v3",
            "license":  "MIT",
		})
	})
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "API is up and running!")
	})
	router.GET("/version", func(c *gin.Context) {
		c.String(200, "API version 1.0")
	})
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	})
}

// func UserRoutes(router *gin.Engine, userController users.UserController){
// 	users := router.Group("/users")
// 	users.POST("/", userController.RegisterUser)
// }
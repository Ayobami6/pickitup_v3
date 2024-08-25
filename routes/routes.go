package routes

import (
	"fmt"
	"net/http"

	_ "github.com/Ayobami6/pickitup_v3/cmd/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// func UserRoutes(router *gin.Engine, userController users.UserController){
// 	users := router.Group("/users")
// 	users.POST("/", userController.RegisterUser)
// }
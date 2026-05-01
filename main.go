package main

import (
	"shop/config"
	"shop/handlers"
	"shop/middleware" // Не забудь создать этот пакет
	"shop/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	err := config.DB.AutoMigrate(
		&models.Brand{},
		&models.Category{},
		&models.Product{},
		&models.User{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProduct)
	r.GET("/brands", handlers.GetBrands)
	r.GET("/categories", handlers.GetCategories)

	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.POST("/products", handlers.CreateProduct)
		auth.PUT("/products/:id", handlers.UpdateProduct)
		auth.DELETE("/products/:id", handlers.DeleteProduct)

		auth.POST("/brands", handlers.CreateBrand)
		auth.POST("/categories", handlers.CreateCategory)
	}

	r.Run(":8080")
}

package main

import (
	"shop/config"
	"shop/handlers"
	"shop/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	if err := config.DB.AutoMigrate(&models.Brand{}, &models.Category{}, &models.Product{}); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	r.GET("/products", handlers.GetProducts)
	r.POST("/products", handlers.CreateProduct)
	r.GET("/products/:id", handlers.GetProduct)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)

	r.GET("/brands", handlers.GetBrands)
	r.POST("/brands", handlers.CreateBrand)

	r.GET("/categories", handlers.GetCategories)
	r.POST("/categories", handlers.CreateCategory)

	r.Run(":8080")
}

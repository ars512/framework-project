package main

import (
	"shop/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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

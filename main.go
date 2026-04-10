package main

import (
	"shop/config"
	"shop/handlers"
	"shop/middleware" // Не забудь создать этот пакет
	"shop/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Инициализация БД
	config.ConnectDB()

	// 2. Миграции (добавили models.User{})
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

	// --- Публичные роуты ---
	r.GET("/health", handlers.HealthCheck)

	// Аутентификация
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Просмотр товаров, брендов и категорий доступен всем
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProduct)
	r.GET("/brands", handlers.GetBrands)
	r.GET("/categories", handlers.GetCategories)

	// --- Защищенные роуты ---
	// Все роуты внутри этой группы требуют валидный JWT токен
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// Управление товарами
		auth.POST("/products", handlers.CreateProduct)
		auth.PUT("/products/:id", handlers.UpdateProduct)
		auth.DELETE("/products/:id", handlers.DeleteProduct)

		// Управление справочниками
		auth.POST("/brands", handlers.CreateBrand)
		auth.POST("/categories", handlers.CreateCategory)
	}

	// Запуск сервера
	r.Run(":8080")
}

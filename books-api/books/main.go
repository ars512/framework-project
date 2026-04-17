package main

import (
	"os"

	"bookstore/handlers"
	"bookstore/middleware"
	"bookstore/store"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	appStore := store.NewStore()
	appStore.SeedDemoData()

	bookHandler := handlers.NewBookHandler(appStore)
	authorHandler := handlers.NewAuthorHandler(appStore)
	categoryHandler := handlers.NewCategoryHandler(appStore)
	favoriteHandler := handlers.NewFavoriteHandler(appStore)

	jwtSecret := getEnv("JWT_SECRET", "secret123")
	authMiddleware := middleware.JWTAuth(jwtSecret)

	// Public routes
	r.GET("/books", bookHandler.ListBooks)
	r.GET("/books/:id", bookHandler.GetBook)

	r.GET("/authors", authorHandler.ListAuthors)
	r.POST("/authors", authorHandler.CreateAuthor)

	r.GET("/categories", categoryHandler.ListCategories)
	r.POST("/categories", categoryHandler.CreateCategory)

	// Book CRUD
	r.POST("/books", bookHandler.CreateBook)
	r.PUT("/books/:id", bookHandler.UpdateBook)
	r.DELETE("/books/:id", bookHandler.DeleteBook)

	// Favorites (JWT protected)
	favorites := r.Group("/books", authMiddleware)
	{
		favorites.GET("/favorites", favoriteHandler.GetFavorites)
		favorites.PUT("/:bookId/favorites", favoriteHandler.AddFavorite)
		favorites.DELETE("/:bookId/favorites", favoriteHandler.RemoveFavorite)
	}

	_ = r.Run(":8080")
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
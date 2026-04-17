package handlers

import (
	"net/http"
	"strconv"

	"bookstore/store"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	Store *store.Store
}

func NewFavoriteHandler(s *store.Store) *FavoriteHandler {
	return &FavoriteHandler{Store: s}
}

func getUserIDFromContext(c *gin.Context) (int, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return 0, false
	}

	switch v := raw.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type"})
		return 0, false
	}
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	page, limit := paginationFromQuery(c)

	items, total := h.Store.ListFavorites(userID, page, limit)

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	_, err = h.Store.AddFavorite(userID, bookID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		case store.ErrInvalidData:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid favorite data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add favorite"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book added to favorites"})
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	removed := h.Store.RemoveFavorite(userID, bookID)
	if !removed {
		c.JSON(http.StatusNotFound, gin.H{"error": "favorite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book removed from favorites"})
}
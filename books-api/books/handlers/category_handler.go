package handlers

import (
	"net/http"

	"bookstore/models"
	"bookstore/store"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Store *store.Store
}

func NewCategoryHandler(s *store.Store) *CategoryHandler {
	return &CategoryHandler{Store: s}
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": h.Store.ListCategories(),
	})
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	category, err := h.Store.CreateCategory(req.Name)
	if err != nil {
		switch err {
		case store.ErrDuplicateName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "category already exists"})
		case store.ErrInvalidData:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		}
		return
	}

	c.JSON(http.StatusCreated, category)
}
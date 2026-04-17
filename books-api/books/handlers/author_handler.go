package handlers

import (
	"net/http"

	"bookstore/models"
	"bookstore/store"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	Store *store.Store
}

func NewAuthorHandler(s *store.Store) *AuthorHandler {
	return &AuthorHandler{Store: s}
}

func (h *AuthorHandler) ListAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": h.Store.ListAuthors(),
	})
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var req models.AuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	author, err := h.Store.CreateAuthor(req.Name)
	if err != nil {
		switch err {
		case store.ErrDuplicateName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "author already exists"})
		case store.ErrInvalidData:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create author"})
		}
		return
	}

	c.JSON(http.StatusCreated, author)
}
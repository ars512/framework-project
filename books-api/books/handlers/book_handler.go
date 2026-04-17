package handlers

import (
	"net/http"
	"strconv"

	"bookstore/models"
	"bookstore/store"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	Store *store.Store
}

func NewBookHandler(s *store.Store) *BookHandler {
	return &BookHandler{Store: s}
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	page, limit := paginationFromQuery(c)

	minPrice := 0.0
	maxPrice := 0.0

	if v := c.Query("min_price"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			minPrice = parsed
		}
	}

	if v := c.Query("max_price"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			maxPrice = parsed
		}
	}

	filter := models.BookFilter{
		Title:    c.Query("title"),
		Author:   c.Query("author"),
		Category: c.Query("category"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Page:     page,
		Limit:    limit,
	}

	items, total := h.Store.ListBooks(filter)

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	book, found := h.Store.GetBookView(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req models.BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if req.Title == "" || req.AuthorID <= 0 || req.CategoryID <= 0 || req.Price < 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, author_id, category_id and valid price are required"})
		return
	}

	book, err := h.Store.CreateBook(req)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "author or category not found"})
		case store.ErrInvalidData:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		}
		return
	}

	view, _ := h.Store.GetBookView(book.ID)
	c.JSON(http.StatusCreated, view)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req models.BookUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if req.Title == "" || req.AuthorID <= 0 || req.CategoryID <= 0 || req.Price < 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, author_id, category_id and valid price are required"})
		return
	}

	book, err := h.Store.UpdateBook(id, req)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "book, author or category not found"})
		case store.ErrInvalidData:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book data"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		}
		return
	}

	view, _ := h.Store.GetBookView(book.ID)
	c.JSON(http.StatusOK, view)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	deleted := h.Store.DeleteBook(id)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}
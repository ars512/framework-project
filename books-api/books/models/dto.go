package models

import "time"

type BookCreateRequest struct {
	Title      string  `json:"title" binding:"required"`
	AuthorID   int     `json:"author_id" binding:"required"`
	CategoryID int     `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
}

type BookUpdateRequest struct {
	Title      string  `json:"title" binding:"required"`
	AuthorID   int     `json:"author_id" binding:"required"`
	CategoryID int     `json:"category_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
}

type AuthorRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type BookFilter struct {
	Title    string
	Author   string
	Category string
	MinPrice float64
	MaxPrice float64
	Page     int
	Limit    int
}

type BookView struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	AuthorID     int     `json:"author_id"`
	AuthorName   string  `json:"author_name"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
}

type FavoriteBookView struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	AuthorID     int       `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
}
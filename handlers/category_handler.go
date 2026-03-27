package handlers

import (
	"net/http"

	"shop/models"

	"github.com/gin-gonic/gin"
)

var categories []models.Category
var categoryID = 1

func GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	category.ID = categoryID
	categoryID++
	categories = append(categories, category)

	c.JSON(http.StatusCreated, category)
}

package handlers

import (
	"net/http"
	"strings"

	"shop/config"
	"shop/models"

	"github.com/gin-gonic/gin"
)

type brandRequest struct {
	Name string `json:"name"`
}

func GetBrands(c *gin.Context) {
	var brands []models.Brand

	if err := config.DB.Order("id asc").Find(&brands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load brands"})
		return
	}

	c.JSON(http.StatusOK, brands)
}

func CreateBrand(c *gin.Context) {
	var req brandRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	brand := models.Brand{Name: name}

	if err := config.DB.Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create brand"})
		return
	}

	c.JSON(http.StatusCreated, brand)
}

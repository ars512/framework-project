package handlers

import (
	"net/http"

	"shop/models"

	"github.com/gin-gonic/gin"
)

var brands []models.Brand
var brandID = 1

func GetBrands(c *gin.Context) {
	c.JSON(http.StatusOK, brands)
}

func CreateBrand(c *gin.Context) {
	var brand models.Brand

	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if brand.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	brand.ID = brandID
	brandID++
	brands = append(brands, brand)

	c.JSON(http.StatusCreated, brand)
}

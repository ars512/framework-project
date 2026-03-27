package handlers

import (
	"net/http"
	"strconv"

	"shop/models"

	"github.com/gin-gonic/gin"
)

var products []models.Product
var productID = 1

func GetProducts(c *gin.Context) {
	category := c.Query("category")
	pageParam := c.Query("page")

	page := 1
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 5
	start := (page - 1) * pageSize
	end := start + pageSize

	filtered := products
	if category != "" {
		var tmp []models.Product
		for _, p := range products {
			if strconv.Itoa(p.CategoryID) == category {
				tmp = append(tmp, p)
			}
		}
		filtered = tmp
	}

	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	c.JSON(http.StatusOK, filtered[start:end])
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if product.Name == "" || product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	product.ID = productID
	productID++
	products = append(products, product)

	c.JSON(http.StatusCreated, product)
}

func GetProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updated models.Product

	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, p := range products {
		if p.ID == id {
			updated.ID = id
			products[i] = updated
			c.JSON(http.StatusOK, updated)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

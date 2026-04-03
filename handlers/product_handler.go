package handlers

import (
	"net/http"
	"strings"

	"shop/config"
	"shop/models"

	"github.com/gin-gonic/gin"
)

type productRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	BrandID     uint    `json:"brand_id"`
	CategoryID  uint    `json:"category_id"`
}

func GetProducts(c *gin.Context) {
	page, err := parsePositiveInt(c.Query("page"), 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := parsePositiveInt(c.Query("limit"), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	if limit > 100 {
		limit = 100
	}

	db := config.DB.Model(&models.Product{}).Preload("Brand").Preload("Category")

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		db = db.Where("name ILIKE ?", "%"+search+"%")
	}

	if value := c.Query("brand_id"); value != "" {
		brandID, err := parseUint(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand_id"})
			return
		}
		db = db.Where("brand_id = ?", brandID)
	}

	if value := c.Query("category_id"); value != "" {
		categoryID, err := parseUint(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		db = db.Where("category_id = ?", categoryID)
	}

	if value := c.Query("min_price"); value != "" {
		minPrice, err := parseFloat(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price"})
			return
		}
		db = db.Where("price >= ?", minPrice)
	}

	if value := c.Query("max_price"); value != "" {
		maxPrice, err := parseFloat(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price"})
			return
		}
		db = db.Where("price <= ?", maxPrice)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count products"})
		return
	}

	var products []models.Product
	if err := db.Order("id desc").Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  products,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func CreateProduct(c *gin.Context) {
	var req productRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Description) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and description are required"})
		return
	}

	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	if req.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stock cannot be negative"})
		return
	}

	if req.BrandID == 0 || req.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "brand_id and category_id are required"})
		return
	}

	var brand models.Brand
	if err := config.DB.First(&brand, req.BrandID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "brand not found"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
		return
	}

	product := models.Product{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		Stock:       req.Stock,
		BrandID:     req.BrandID,
		CategoryID:  req.CategoryID,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	if err := config.DB.Preload("Brand").Preload("Category").First(&product, product.ID).Error; err != nil {
		c.JSON(http.StatusCreated, product)
		return
	}

	c.JSON(http.StatusCreated, product)
}

func GetProduct(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var product models.Product
	if err := config.DB.Preload("Brand").Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var existing models.Product
	if err := config.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Description) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and description are required"})
		return
	}

	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	if req.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stock cannot be negative"})
		return
	}

	if req.BrandID == 0 || req.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "brand_id and category_id are required"})
		return
	}

	var brand models.Brand
	if err := config.DB.First(&brand, req.BrandID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "brand not found"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
		return
	}

	existing.Name = strings.TrimSpace(req.Name)
	existing.Description = strings.TrimSpace(req.Description)
	existing.Price = req.Price
	existing.Stock = req.Stock
	existing.BrandID = req.BrandID
	existing.CategoryID = req.CategoryID

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	if err := config.DB.Preload("Brand").Preload("Category").First(&existing, id).Error; err != nil {
		c.JSON(http.StatusOK, existing)
		return
	}

	c.JSON(http.StatusOK, existing)
}

func DeleteProduct(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

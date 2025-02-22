package controllers

import (
	"net/http"
	"strconv"

	"bonus/models" // Ganti "bonus" sesuai module name Anda
	"github.com/gin-gonic/gin"
)

// KpiCategoryInput adalah payload untuk input pembuatan / update kategori KPI
type KpiCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateKpiCategory - POST /api/kpi-categories
func CreateKpiCategory(c *gin.Context) {
	var input KpiCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.KpiCategory{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori KPI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// GetKpiCategories - GET /api/kpi-categories
func GetKpiCategories(c *gin.Context) {
	var categories []models.KpiCategory
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori KPI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// UpdateKpiCategory - PUT /api/kpi-categories/:id
func UpdateKpiCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.KpiCategory

	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori KPI tidak ditemukan"})
		return
	}

	var input KpiCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = input.Name
	category.Description = input.Description

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori KPI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteKpiCategory - DELETE /api/kpi-categories/:id
func DeleteKpiCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var category models.KpiCategory

	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori KPI tidak ditemukan"})
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori KPI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

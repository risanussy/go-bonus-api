package controllers

import (
	"net/http"
	"strconv"

	"bonus/models"

	"github.com/gin-gonic/gin"
)

// GET /api/kpis
func GetKPIs(c *gin.Context) {
	var kpis []models.KPI
	if err := db.Find(&kpis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data KPI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kpis})
}

// POST /api/kpis
func CreateKPI(c *gin.Context) {
	var input models.KPI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	// Set EmployeeID dari context JWT (jika diperlukan)
	if empID, exists := c.Get("employee_id"); exists {
		// Pastikan konversi tipe sesuai dengan tipe data EmployeeID
		input.EmployeeID = uint(empID.(float64))
	}
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat KPI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// PUT /api/kpis/:id
func UpdateKPI(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var kpi models.KPI
	if err := db.First(&kpi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KPI tidak ditemukan"})
		return
	}
	var input models.KPI
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}
	kpi.Title = input.Title
	kpi.Description = input.Description
	kpi.Score = input.Score
	kpi.Validated = input.Validated
	db.Save(&kpi)
	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// DELETE /api/kpis/:id
func DeleteKPI(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var kpi models.KPI
	if err := db.First(&kpi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KPI tidak ditemukan"})
		return
	}
	db.Delete(&kpi)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

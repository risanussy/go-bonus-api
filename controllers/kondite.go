package controllers

import (
	"net/http"
	"time"

	"bonus/models"

	"github.com/gin-gonic/gin"
)

// GetKondites - GET /api/kondites
// Mengambil daftar kondite (opsional: boleh filter by employee_id).
// GET /api/kondites
func GetKondites(c *gin.Context) {
	var kondites []models.Kondite
	if err := db.Preload("Employee").Find(&kondites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kondite"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kondites})
}

// CreateKondite - POST /api/kondites
// Membuat kondite baru (SP1, SP2, dsb.).
func CreateKondite(c *gin.Context) {
	var input struct {
		EmployeeID  uint   `json:"employee_id" binding:"required"`
		Category    string `json:"category" binding:"required"`
		StartDate   string `json:"start_date" binding:"required"` // Format "YYYY-MM-DD"
		EndDate     string `json:"end_date" binding:"required"`   // Format "YYYY-MM-DD"
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse tanggal
	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format start_date tidak valid (YYYY-MM-DD)"})
		return
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format end_date tidak valid (YYYY-MM-DD)"})
		return
	}

	kondite := models.Kondite{
		EmployeeID:  input.EmployeeID,
		Category:    input.Category,
		StartDate:   start,
		EndDate:     end,
		Description: input.Description,
	}

	if err := db.Create(&kondite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kondite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kondite})
}

// UpdateKondite - PUT /api/kondites/:id
// Memperbarui data kondite tertentu.
func UpdateKondite(c *gin.Context) {
	id := c.Param("id")

	// Cari data existing
	var kondite models.Kondite
	if err := db.First(&kondite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kondite tidak ditemukan"})
		return
	}

	var input struct {
		EmployeeID  uint   `json:"employee_id"`
		Category    string `json:"category"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Jika field kosong, biarkan data lama
	if input.EmployeeID != 0 {
		kondite.EmployeeID = input.EmployeeID
	}
	if input.Category != "" {
		kondite.Category = input.Category
	}
	if input.StartDate != "" {
		start, err := time.Parse("2006-01-02", input.StartDate)
		if err == nil {
			kondite.StartDate = start
		}
	}
	if input.EndDate != "" {
		end, err := time.Parse("2006-01-02", input.EndDate)
		if err == nil {
			kondite.EndDate = end
		}
	}
	if input.Description != "" {
		kondite.Description = input.Description
	}

	if err := db.Save(&kondite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data kondite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kondite})
}

// DeleteKondite - DELETE /api/kondites/:id
// Menghapus data kondite.
func DeleteKondite(c *gin.Context) {
	id := c.Param("id")

	var kondite models.Kondite
	if err := db.First(&kondite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kondite tidak ditemukan"})
		return
	}

	if err := db.Delete(&kondite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kondite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Kondite berhasil dihapus"})
}

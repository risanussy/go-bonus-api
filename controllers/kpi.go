package controllers

import (
    "net/http"
    "strconv"

    "bonus/models"

    "github.com/gin-gonic/gin"
)

// GET /api/kpis
// Ambil semua KPI
func GetKPIs(c *gin.Context) {
    var kpis []models.KPI
    if err := db.Find(&kpis).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": kpis})
}

// GET /api/kpis/:id
// Ambil satu KPI berdasar ID
func GetKPIByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var kpi models.KPI
    if err := db.First(&kpi, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "KPI tidak ditemukan"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// POST /api/kpis
// Buat KPI baru
func CreateKPI(c *gin.Context) {
    var input models.KPI
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
        return
    }

    if err := db.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": input})
}

// PUT /api/kpis/:id
// Update KPI yang sudah ada
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

    // Update field sesuai input
    kpi.Title       = input.Title
    kpi.Category    = input.Category
    kpi.Weight      = input.Weight
    kpi.Target      = input.Target
    kpi.Poor        = input.Poor
    kpi.Fair        = input.Fair
    kpi.Good        = input.Good
    kpi.Outstanding = input.Outstanding
    kpi.Exceptional = input.Exceptional
    kpi.Score       = input.Score
    kpi.Validated   = input.Validated
    kpi.EmployeeID  = input.EmployeeID

    if err := db.Save(&kpi).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": kpi})
}

// DELETE /api/kpis/:id
// Hapus KPI berdasarkan ID
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

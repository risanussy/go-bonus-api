package controllers

import (
	"net/http"

	"bonus/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AHPRequest merepresentasikan payload untuk perhitungan bonus (jika menggunakan bobot AHP).
type AHPRequest struct {
	CriteriaWeights map[string]float64 `json:"criteria_weights"`
	// Contoh pemakaian: {"KPI": 0.5, "Kondite": 0.3, "PenambahPoin": 0.2}
}

// GetBonus - GET /api/bonus
// Mengambil data bonus setiap pegawai (dihitung langsung).
func GetBonus(c *gin.Context) {
	var employees []models.Employee
	if err := db.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}

	// Hasil final akan kita tampung di "results".
	var results []map[string]interface{}

	for _, emp := range employees {
		// Ambil data KPI milik pegawai ini
		var kpis []models.KPI
		if err := db.Where("employee_id = ?", emp.ID).Find(&kpis).Error; err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		// Contoh sederhana: totalScore adalah penjumlahan semua Score KPI.
		totalScore := 0.0
		for _, kpi := range kpis {
			totalScore += kpi.Score
		}

		// Contoh perhitungan bonus = totalScore * 1000 (silakan ganti rumus sesuai kebutuhan)
		bonus := totalScore * 1000

		// Masukkan data ke "results"
		results = append(results, gin.H{
			"employee_id": emp.ID,
			"name":        emp.Name,
			"total_score": totalScore,
			"bonus":       bonus,
		})
	}

	// Kembalikan data dalam bentuk JSON
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// CalculateBonus - POST /api/bonus/calculate
// Menghitung bonus berdasarkan bobot kriteria (AHP) atau logika lain.
func CalculateBonus(c *gin.Context) {
	var req AHPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Ambil seluruh pegawai
	var employees []models.Employee
	if err := db.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}

	var results []map[string]interface{}

	// Lakukan perhitungan bonus untuk setiap pegawai
	for _, emp := range employees {
		var kpis []models.KPI
		if err := db.Where("employee_id = ?", emp.ID).Find(&kpis).Error; err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		// Contoh perhitungan totalScore = penjumlahan Score KPI
		totalScore := 0.0
		for _, kpi := range kpis {
			totalScore += kpi.Score
		}

		// Jika menggunakan bobot AHP, Anda bisa mengalikan totalScore dengan bobot tertentu:
		// Misalnya: bonus = totalScore * req.CriteriaWeights["KPI"] * KONSTANTA
		// Di sini untuk contoh sederhana, kita abaikan bobot:
		bonus := totalScore * 1000

		// Tambahkan ke results
		results = append(results, gin.H{
			"employee_id": emp.ID,
			"name":        emp.Name,
			"total_score": totalScore,
			"bonus":       bonus,
		})
	}

	// Kembalikan hasil perhitungan bonus
	c.JSON(http.StatusOK, gin.H{"data": results})
}

package controllers

import (
	"net/http"

	"bonus/models"
	"github.com/gin-gonic/gin"
)

// Struktur input (payload) yang diterima dari front-end
type KPIEvaluationInput struct {
	EmployeeID  uint   `json:"employee_id" binding:"required"`
	KPIID       uint   `json:"kpi_id" binding:"required"`
	Achievement string `json:"achievement" binding:"required"`
}

// Buat daftar tingkatan penilaian => "poor 1", "fair 2", dsb.
var AchievementList = []string{
	"poor 1",
	"fair 2",
	"good 3",
	"outstanding 4",
	"exceptional 5",
}

// GetKPIAchievementList - GET /api/kpi_achievement_list
// Mengembalikan daftar tingkatan penilaian untuk front-end
func GetKPIAchievementList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": AchievementList})
}

// CreateKPIEvaluation - POST /api/kpi_evaluations
// Menerima penilaian KPI dari front-end, mengonversi "achievement" jadi "point"
func CreateKPIEvaluation(c *gin.Context) {
	var input KPIEvaluationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tentukan point berdasarkan achievement
	point := parseAchievementToPoint(input.Achievement)

	kpiev := models.KPIEvaluation{
		EmployeeID:  input.EmployeeID,
		KPIID:       input.KPIID,
		Achievement: input.Achievement, // simpan teks aslinya
		Point:       point,            // simpan nilai numeriknya
	}

	if err := db.Create(&kpiev).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan penilaian KPI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kpiev})
}

// parseAchievementToPoint mengubah "poor 1" -> 1, "fair 2" -> 2, dst.
func parseAchievementToPoint(ach string) float64 {
	switch ach {
	case "poor 1":
		return 1
	case "fair 2":
		return 2
	case "good 3":
		return 3
	case "outstanding 4":
		return 4
	case "exceptional 5":
		return 5
	default:
		return 0
	}
}

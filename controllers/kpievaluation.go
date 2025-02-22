package controllers

import (
	"net/http"

	"bonus/models"
	"github.com/gin-gonic/gin"
)

// Struct input (payload) yang diterima dari front-end
type KPIEvaluationInput struct {
	EmployeeID  uint   `json:"employee_id" binding:"required"`
	KPIID       uint   `json:"kpi_id" binding:"required"`
	Achievement string `json:"achievement" binding:"required"`
}

// CreateKPIEvaluation menerima penilaian KPI dari front-end
func CreateKPIEvaluation(c *gin.Context) {
	var input KPIEvaluationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpiev := models.KPIEvaluation{
		EmployeeID:  input.EmployeeID,
		KPIID:       input.KPIID,
		Achievement: input.Achievement,
	}

	if err := db.Create(&kpiev).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan penilaian KPI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kpiev})
}

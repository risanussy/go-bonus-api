// models/kpi_evaluation.go
package models

import "gorm.io/gorm"

type KPIEvaluation struct {
	gorm.Model
	EmployeeID  uint   `json:"employee_id"`
	KPIID       uint   `json:"kpi_id"`
	Achievement string `json:"achievement"`
	Point       float64 `json:"point"` // menampung nilai numeric dari achievement
}

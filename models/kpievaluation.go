package models

import "gorm.io/gorm"

// KPIEvaluation merepresentasikan penilaian KPI tertentu
// oleh seorang karyawan, dengan nilai pencapaian tertentu.
type KPIEvaluation struct {
	gorm.Model
	EmployeeID  uint   `json:"employee_id"`
	KPIID       uint   `json:"kpi_id"`
	Achievement string `json:"achievement"`
}

package models

import "gorm.io/gorm"

// KpiCategory merepresentasikan kategori KPI (mis. "KPI Individu", "KPI Departemen", dsb.)
type KpiCategory struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

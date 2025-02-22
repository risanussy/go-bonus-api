package models

import "gorm.io/gorm"

type KPI struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Score       float64 `json:"score"`       // Nilai KPI yang diinput pegawai
	Validated   bool    `json:"validated"`   // Validasi oleh atasan
	EmployeeID  uint    `json:"employee_id"` // Relasi ke pegawai
}

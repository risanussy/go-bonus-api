package models

import "gorm.io/gorm"

type KPI struct {
	gorm.Model
	Title        string  `json:"title"`
	Category     string  `json:"category"`
	Weight       float64 `json:"weight"`
	Target       string  `json:"target"`
	Poor         string  `json:"poor"`
	Fair         string  `json:"fair"`
	Good         string  `json:"good"`
	Outstanding  string  `json:"outstanding"`
	Exceptional  string  `json:"exceptional"`

	Score      float64 `json:"score"`       // Nilai KPI yang diinput pegawai
	Validated  bool    `json:"validated"`   // Validasi oleh atasan
	EmployeeID uint    `json:"employee_id"` // Relasi ke pegawai
}

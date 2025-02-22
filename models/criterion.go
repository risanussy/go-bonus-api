package models

import "gorm.io/gorm"

// Criterion merepresentasikan kriteria penilaian.
type Criterion struct {
	gorm.Model
	Name        string  `json:"name"`        // Nama kriteria (misalnya: "KPI", "Kondite", "Faktor Lain")
	Description string  `json:"description"` // Deskripsi kriteria
	Weight      float64 `json:"weight"`      // Bobot default yang dapat digunakan pada perhitungan AHP
}

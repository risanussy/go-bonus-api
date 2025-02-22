package models

import "gorm.io/gorm"

// Evaluation menyimpan hasil evaluasi kinerja pegawai per kriteria.
type Evaluation struct {
	gorm.Model
	EmployeeID  uint    `json:"employee_id"`  // ID pegawai yang dievaluasi
	CriterionID uint    `json:"criterion_id"` // ID kriteria yang digunakan dalam evaluasi
	Score       float64 `json:"score"`        // Nilai atau skor yang diberikan
	Comments    string  `json:"comments"`     // Komentar tambahan (opsional)

	// Relasi (opsional) untuk mendapatkan data lengkap pegawai dan kriteria
	Employee  Employee  `gorm:"foreignKey:EmployeeID" json:"employee"`
	Criterion Criterion `gorm:"foreignKey:CriterionID" json:"criterion"`
}

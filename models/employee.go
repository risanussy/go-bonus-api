// models/employee.go
package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"` // Simpan hash password
	Role     string `json:"role"`     // Contoh: admin, HRD, manager, pegawai
	Salary   int    `json:"salary"`   // Tambahkan field gaji
}

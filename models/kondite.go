package models

import (
	"time"

	"gorm.io/gorm"
)

type Kondite struct {
	gorm.Model
	EmployeeID  uint      `json:"employee_id"`
	Category    string    `json:"category"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`

	// Relasi ke Employee
	Employee Employee `json:"employee" gorm:"foreignKey:EmployeeID"`
}

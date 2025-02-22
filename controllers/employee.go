package controllers

import (
	"net/http"
	"bonus/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetEmployees mengembalikan daftar semua pegawai.
func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := db.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": employees})
}

// CreateEmployee membuat data pegawai baru (dengan Salary).
func CreateEmployee(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Salary   int    `json:"salary" binding:"required"` // Field gaji
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	employee := models.Employee{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
		Salary:   input.Salary, // Simpan salary
	}

	if err := db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat data pegawai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employee})
}

// UpdateEmployee memperbarui data pegawai yang ada (termasuk Salary).
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := db.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	var input struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Salary int    `json:"salary"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.Name = input.Name
	employee.Email = input.Email
	employee.Role = input.Role
	employee.Salary = input.Salary // Update gaji

	if err := db.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pegawai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employee})
}

// DeleteEmployee menghapus data pegawai.
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := db.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
		return
	}

	if err := db.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pegawai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Pegawai berhasil dihapus"})
}

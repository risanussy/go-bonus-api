package controllers

import (
	"net/http"
	"time"

	"bonus/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtKey = []byte("your_secret_key") // Ubah sesuai dengan secret key

// SetDB digunakan untuk menginisialisasi instance DB pada controllers
func SetDB(database *gorm.DB) {
	db = database
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	var employee models.Employee
	if err := db.Where("email = ?", input.Email).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Bandingkan password (pastikan password sudah di-hash saat registrasi)
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Buat JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"employee_id": employee.ID,
		"role":        employee.Role,
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

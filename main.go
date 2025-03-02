package main

import (
	"log"
	"net/http"
	"time"

	"bonus/controllers"
	"bonus/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()

	// Middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Koneksi ke DB MySQL
	dsn := "root@tcp(127.0.0.1:3306)/order_bonus_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Migrasi model
	db.AutoMigrate(
		&models.Employee{},
		&models.KPI{},
		&models.Criterion{},
		&models.Evaluation{},
		&models.KPIEvaluation{},
		&models.KpiCategory{},
		&models.Kondite{},
	)

	// Seed data admin setelah migrasi
	seedAdmin(db)

	// Set DB di controllers
	controllers.SetDB(db)

	api := router.Group("/api")
	{
		// Bonus
		api.POST("/bonus/calculate", controllers.CalculateBonus)
		api.GET("/bonus", controllers.GetBonus)

		// Kalibrasi
		api.GET("/kalibrasi", controllers.GetKalibrasi)

		// Login (auth)
		api.POST("/login", controllers.Login)

		// KPI
		api.GET("/kpis", controllers.GetKPIs)
		api.POST("/kpis", controllers.CreateKPI)
		api.PUT("/kpis/:id", controllers.UpdateKPI)
		api.DELETE("/kpis/:id", controllers.DeleteKPI)

		// Kategori KPI
		api.GET("/kpi-categories", controllers.GetKpiCategories)
		api.POST("/kpi-categories", controllers.CreateKpiCategory)
		api.PUT("/kpi-categories/:id", controllers.UpdateKpiCategory)
		api.DELETE("/kpi-categories/:id", controllers.DeleteKpiCategory)

		// Penilaian KPI
		api.GET("/kpi_evaluations", controllers.GetKPIAchievementList)
		api.POST("/kpi_evaluations", controllers.CreateKPIEvaluation)

		// Employee
		api.GET("/employees", controllers.GetEmployees)
		api.POST("/employees", controllers.CreateEmployee)
		api.PUT("/employees/:id", controllers.UpdateEmployee)
		api.DELETE("/employees/:id", controllers.DeleteEmployee)
		// Kondite
		api.GET("/kondites", controllers.GetKondites)
		api.POST("/kondites", controllers.CreateKondite)
		api.PUT("/kondites/:id", controllers.UpdateKondite)
		api.DELETE("/kondites/:id", controllers.DeleteKondite)
	}

	// Jalankan server di port 8080
	http.ListenAndServe(":8080", router)
}

// seedAdmin membuat data admin default jika belum ada
func seedAdmin(db *gorm.DB) {
	// Cek apakah ada admin (email = "admin@admin.com")
	var existing models.Employee
	if err := db.Where("email = ?", "admin@admin.com").First(&existing).Error; err != nil {
		// Jika not found, buat admin
		if err == gorm.ErrRecordNotFound {
			// Hash password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1234"), bcrypt.DefaultCost)

			admin := models.Employee{
				Name:     "admin",
				Email:    "admin@admin.com",
				Password: string(hashedPassword),
				Role:     "admin",
				Salary:   4000000, // gaji 4.000.000
			}
			db.Create(&admin)
			log.Println("Admin default berhasil dibuat (admin@admin.com / admin123)")
		} else {
			log.Println("Gagal cek data admin:", err)
		}
	} else {
		log.Println("Admin default sudah ada, tidak perlu seed ulang.")
	}
}

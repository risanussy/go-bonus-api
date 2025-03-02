package controllers

import (
	"math"
	"net/http"

	"bonus/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KalibrasiResponse merepresentasikan struktur data yang akan dikembalikan ke front-end.
type KalibrasiResponse struct {
	No                  int     `json:"no"`
	Name                string  `json:"name"`
	KPIPerusahaan       float64 `json:"kpi_perusahaan"`
	KPIDepart           float64 `json:"kpi_depart"`
	KPIIndividu         float64 `json:"kpi_individu"`
	TotalKPI            float64 `json:"total_kpi"`
	PengurangPoin       float64 `json:"pengurang_poin"`
	PenambahPoin        float64 `json:"penambah_poin"`
	KPISetelahKalibrasi float64 `json:"kpi_setelah_kalibrasi"`
	Skala               string  `json:"skala"`
	Gaji                float64 `json:"gaji"`
	Bonus               float64 `json:"bonus"`
}

// GetKalibrasi - GET /api/kalibrasi
// Menampilkan hasil kalibrasi KPI setiap karyawan.
// - Admin melihat semua karyawan.
// - User biasa hanya melihat miliknya sendiri.
func GetKalibrasi(c *gin.Context) {
	// Ambil role & employee_id dari middleware JWT
	roleVal, roleExists := c.Get("role")
	empVal, empExists := c.Get("employee_id")
	if !roleExists || !empExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := roleVal.(string)
	currentUserID := empVal.(uint)

	// Ambil data employees
	var employees []models.Employee
	if role == "admin" {
		// Admin => semua karyawan
		if err := db.Find(&employees).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
			return
		}
	} else {
		// User => hanya data karyawan untuk dirinya sendiri
		var emp models.Employee
		if err := db.First(&emp, currentUserID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pegawai tidak ditemukan"})
			return
		}
		employees = []models.Employee{emp}
	}

	results := []KalibrasiResponse{}
	nomor := 1

	// Loop setiap pegawai => hitung KPI & bonus
	for _, emp := range employees {
		// Ambil KPI milik pegawai ini
		var kpis []models.KPI
		err := db.Where("employee_id = ?", emp.ID).Find(&kpis).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		// Kita pisahkan total KPI berdasarkan kategori
		var totalPerusahaan float64
		var totalDept float64
		var totalInd float64

		// Loop setiap KPI, hitung finalScore = Score * (Weight / 100)
		for _, k := range kpis {
			finalScore := k.Score * (k.Weight / 100.0)

			switch k.Category {
			case "Perusahaan":
				totalPerusahaan += finalScore
			case "Dept", "Departemen":
				totalDept += finalScore
			case "Individu":
				totalInd += finalScore
			}
		}

		// Total KPI sebelum kalibrasi
		totalKPI := totalPerusahaan + totalDept + totalInd

		// Pengurang poin => dari Kondite
		pengurang, _ := HitungPengurangPoin(emp.ID)

		// Penambah poin => jika ada reward dsb. (default 0)
		penambah, _ := HitungPenambahPoin(emp.ID)

		// Final KPI setelah kalibrasi
		finalKPI := totalKPI - pengurang + penambah
		if finalKPI < 0 {
			finalKPI = 0
		}

		// Skala & multiplier bonus
		skala, multiplier := SkalaKPI(finalKPI)

		// Gaji
		gaji := float64(emp.Salary)

		// Bonus
		bonus := gaji * multiplier

		// Buat item response
		item := KalibrasiResponse{
			No:                  nomor,
			Name:                emp.Name,
			KPIPerusahaan:       RoundFloat(totalPerusahaan, 1),
			KPIDepart:           RoundFloat(totalDept, 1),
			KPIIndividu:         RoundFloat(totalInd, 1),
			TotalKPI:            RoundFloat(totalKPI, 1),
			PengurangPoin:       RoundFloat(pengurang, 1),
			PenambahPoin:        RoundFloat(penambah, 1),
			KPISetelahKalibrasi: RoundFloat(finalKPI, 1),
			Skala:               skala,
			Gaji:                gaji,
			Bonus:               bonus,
		}
		results = append(results, item)
		nomor++
	}

	// Return JSON
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// HitungPengurangPoin => contoh perhitungan total min_point dari Kondite
func HitungPengurangPoin(empID uint) (float64, error) {
	var kondites []models.Kondite
	err := db.Where("employee_id = ?", empID).Find(&kondites).Error
	if err != nil {
		return 0, err
	}
	var total float64
	for _, k := range kondites {
		// Asumsi MinPoint adalah float64 atau int => sesuaikan
		total += float64(k.MinPoint)
	}
	return total, nil
}

// HitungPenambahPoin => contoh, jika ada reward dsb.
func HitungPenambahPoin(empID uint) (float64, error) {
	// Default 0, diisi jika ada logika penambahan
	return 0, nil
}

// SkalaKPI => mengembalikan keterangan & multiplier bonus
// < 2 => "Poor" => multiplier=1
// < 3 => "Fair" => multiplier=2
// < 4 => "Good" => multiplier=3
// < 5 => "Outstanding" => multiplier=4
// >= 5 => "Exceptional" => multiplier=5
func SkalaKPI(final float64) (string, float64) {
	if final < 2 {
		return "Poor", 1
	} else if final < 3 {
		return "Fair", 2
	} else if final < 4 {
		return "Good", 3
	} else if final < 5 {
		return "Outstanding", 4
	}
	return "Exceptional", 5
}

// RoundFloat => membulatkan float ke n desimal
func RoundFloat(val float64, places int) float64 {
	if places < 0 {
		return val
	}
	shift := math.Pow(10, float64(places))
	return math.Round(val*shift) / shift
}

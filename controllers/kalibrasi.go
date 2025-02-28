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
// Menghitung nilai KPI (Perusahaan, Dept, Individu), pengurang poin, bonus, dll.
func GetKalibrasi(c *gin.Context) {
	// 1. Ambil semua pegawai
	var employees []models.Employee
	if err := db.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pegawai"})
		return
	}

	results := []KalibrasiResponse{}
	nomor := 1

	// 2. Loop setiap pegawai => hitung KPI & bonus
	for _, emp := range employees {
		// 2a. Ambil KPI milik pegawai ini
		// Pastikan di struct `KPI` ada field `Category` atau cara lain
		// untuk memisahkan "Perusahaan", "Dept", "Individu"
		var kpis []models.KPI
		err := db.Where("employee_id = ?", emp.ID).Find(&kpis).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		// 2b. Bagi KPI berdasarkan category/description
		var totalPerusahaan float64
		var totalDept float64
		var totalInd float64

		for _, k := range kpis {
			// Misal k.Description: "Perusahaan", "Dept", "Individu"
			switch k.Description {
			case "Perusahaan":
				totalPerusahaan += k.Score
			case "Dept", "Departemen":
				totalDept += k.Score
			case "Individu":
				totalInd += k.Score
			}
		}

		// 2c. Total KPI (sebelum kalibrasi)
		totalKPI := totalPerusahaan + totalDept + totalInd

		// 2d. Hitung pengurang poin => dari Kondite
		pengurang, _ := HitungPengurangPoin(emp.ID)

		// 2e. Hitung penambah poin (opsional, misal reward)
		penambah, _ := HitungPenambahPoin(emp.ID)

		// 2f. Final KPI
		finalKPI := totalKPI - pengurang + penambah
		if finalKPI < 0 {
			finalKPI = 0
		}

		// 2g. Skala (Poor, Fair, Good, Outstanding, Exceptional)
		skala, multiplier := SkalaKPI(finalKPI)

		// 2h. Gaji
		gaji := float64(emp.Salary)

		// 2i. Bonus = gaji * multiplier
		bonus := gaji * multiplier

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

	// 3. Return JSON
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
		total += float64(k.MinPoint) // Asumsi MinPoint = int64
	}
	return total, nil
}

// HitungPenambahPoin => contoh, jika ada reward dsb.
func HitungPenambahPoin(empID uint) (float64, error) {
	// Di sini diisi logika penambah poin, misalnya "BOD Award" = 0.5
	// Contoh default: 0
	return 0, nil
}

// SkalaKPI => mengembalikan keterangan & multiplier
// Contoh rumus:
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

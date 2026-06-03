package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import driver postgres
)

// Konfigurasi Database
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres" 
	password = "admin7841" 
	dbname   = "bioskop_db"
)

// Struct Bioskop dengan tag json dan binding validasi dari Gin
type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama" binding:"required"`   // binding:"required" = tidak boleh kosong
	Lokasi string  `json:"lokasi" binding:"required"` // binding:"required" = tidak boleh kosong
	Rating float64 `json:"rating"`
}

func main() {
	// =========================================================
	// 1. KONEKSI KE DATABASE POSTGRESQL
	// =========================================================
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Memastikan koneksi benar-benar terhubung
	err = db.Ping()
	if err != nil {
		panic("Gagal terhubung ke database! Cek password atau nama database: " + err.Error())
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	// =========================================================
	// 2. INISIALISASI ROUTER GIN
	// =========================================================
	router := gin.Default()

	// =========================================================
	// 3. ENDPOINT POST /bioskop
	// =========================================================
	router.POST("/bioskop", func(c *gin.Context) {
		var b Bioskop

		// ShouldBindJSON otomatis mengecek apakah "nama" dan "lokasi" ada isinya
		if err := c.ShouldBindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Validasi Gagal: Nama dan Lokasi tidak boleh kosong!",
			})
			return
		}

		// Query SQL untuk memasukkan data dan mengembalikan ID yang baru dibuat
		sqlStatement := `
		INSERT INTO bioskop (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		RETURNING id`

		id := 0
		// Mengeksekusi query dan menyimpan ID baru ke variabel 'id'
		err = db.QueryRow(sqlStatement, b.Nama, b.Lokasi, b.Rating).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal menyimpan ke database: " + err.Error(),
			})
			return
		}

		b.ID = id
		c.JSON(http.StatusCreated, gin.H{
			"message": "Data Bioskop berhasil ditambahkan!",
			"data":    b,
		})
	})

	// Menjalankan server Gin di port 8080
	fmt.Println("Server Gin berjalan di http://localhost:8080")
	router.Run(":8080")
}
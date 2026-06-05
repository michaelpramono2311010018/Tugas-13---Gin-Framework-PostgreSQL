package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Konfigurasi Database
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin7841"
	dbname   = "bioskop_db"
)

type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama" binding:"required"`
	Lokasi string  `json:"lokasi" binding:"required"`
	Rating float64 `json:"rating"`
}

func main() {
	// =========================================================
	// 1. KONEKSI DATABASE
	// =========================================================
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic("Gagal terhubung ke database! Cek password: " + err.Error())
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	router := gin.Default()

	// =========================================================
	// ENDPOINT: CREATE (POST) - Dari Tugas 13
	// =========================================================
	router.POST("/bioskop", func(c *gin.Context) {
		var b Bioskop
		if err := c.ShouldBindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi Gagal: Nama dan Lokasi tidak boleh kosong!"})
			return
		}

		sqlStatement := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id`
		id := 0
		err = db.QueryRow(sqlStatement, b.Nama, b.Lokasi, b.Rating).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data: " + err.Error()})
			return
		}

		b.ID = id
		c.JSON(http.StatusCreated, gin.H{
			"message": "Data Bioskop berhasil ditambahkan!",
			"data":    b,
		})
	})

	// =========================================================
	// ENDPOINT: READ ALL (GET /bioskop)
	// =========================================================
	router.GET("/bioskop", func(c *gin.Context) {
		var results []Bioskop

		sqlStatement := `SELECT id, nama, lokasi, rating FROM bioskop`
		rows, err := db.Query(sqlStatement)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data: " + err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var b Bioskop
			err = rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca baris data: " + err.Error()})
				return
			}
			results = append(results, b)
		}

		if results == nil {
			results = []Bioskop{}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Berhasil mengambil semua data bioskop",
			"data":    results,
		})
	})

	// =========================================================
	// ENDPOINT: READ BY ID (GET /bioskop/:id)
	// =========================================================
	router.GET("/bioskop/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka valid"})
			return
		}

		var b Bioskop
		sqlStatement := `SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1`
		err = db.QueryRow(sqlStatement, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Data bioskop tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan pada server: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Berhasil mengambil detail bioskop",
			"data":    b,
		})
	})

	// =========================================================
	// ENDPOINT: UPDATE (PUT /bioskop/:id)
	// =========================================================
	router.PUT("/bioskop/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka valid"})
			return
		}

		var b Bioskop
		// Validasi input JSON
		if err := c.ShouldBindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validasi Gagal: Nama dan Lokasi tidak boleh kosong!"})
			return
		}

		sqlStatement := `UPDATE bioskop SET nama = $2, lokasi = $3, rating = $4 WHERE id = $1`
		res, err := db.Exec(sqlStatement, id, b.Nama, b.Lokasi, b.Rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data: " + err.Error()})
			return
		}

		count, err := res.RowsAffected()
		if err != nil || count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data bioskop tidak ditemukan atau tidak ada perubahan"})
			return
		}

		b.ID = id
		c.JSON(http.StatusOK, gin.H{
			"message": "Data bioskop berhasil diperbarui",
			"data":    b,
		})
	})

	// =========================================================
	// ENDPOINT: DELETE (DELETE /bioskop/:id)
	// =========================================================
	router.DELETE("/bioskop/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka valid"})
			return
		}

		sqlStatement := `DELETE FROM bioskop WHERE id = $1`
		res, err := db.Exec(sqlStatement, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data: " + err.Error()})
			return
		}

		count, err := res.RowsAffected()
		if err != nil || count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data bioskop tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Data bioskop berhasil dihapus"})
	})

	fmt.Println("Server Gin berjalan di http://localhost:8080")
	router.Run(":8080")
}
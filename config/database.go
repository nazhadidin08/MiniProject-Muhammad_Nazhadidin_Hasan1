package config

import (
	"ecommerce-api/models"
	"fmt"
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB menginisialisasi koneksi database dan migrasi
func ConnectDB() {
	var err error
	dsn := os.Getenv("DB_DSN")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection successfully opened")

	// (Soal 1-8) Menjalankan AutoMigrate untuk semua model
	err = DB.AutoMigrate(
		&models.User{},
		&models.Toko{},
		&models.Alamat{},
		&models.Kategori{},
		&models.Produk{},
		&models.Transaksi{},
		&models.DetailTransaksi{},
		&models.LogProduk{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated")
}
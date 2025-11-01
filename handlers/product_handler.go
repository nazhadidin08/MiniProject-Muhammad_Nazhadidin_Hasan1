package handlers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"strconv"
)

// (Soal 7) Create Product
func CreateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// (Ketentuan 14) Dapatkan Toko milik user yang login
	var toko models.Toko
	if err := config.DB.First(&toko, "user_id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User does not own a shop"})
	}

	// (Ketentuan 5) Parse multipart/form-data
	name := c.FormValue("name")
	description := c.FormValue("description")
	priceStr := c.FormValue("price")
	stockStr := c.FormValue("stock")
	kategoriIDStr := c.FormValue("kategori_id")

	// Konversi string ke tipe data yang benar
	price, _ := strconv.ParseFloat(priceStr, 64)
	stock, _ := strconv.Atoi(stockStr)
	kategoriID, _ := strconv.ParseUint(kategoriIDStr, 10, 64)

	// (Ketentuan 5) Handle File Upload
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image upload failed"})
	}

	// Buat nama file unik (misal: product-123456.jpg)
	filename := fmt.Sprintf("product-%d%s",
		toko.ID,
		filepath.Ext(file.Filename))
	
	// Simpan file ke folder ./uploads/
	// Pastikan folder 'uploads' ada
	filePath := filepath.Join("uploads", filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	product := models.Produk{
		TokoID:      toko.ID, // (Ketentuan 14)
		KategoriID:  uint(kategoriID),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImagePath:   filePath, // Simpan path-nya
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}



	return c.Status(fiber.StatusCreated).JSON(product)
}

// (Soal 9 & 10) Get All Products with Pagination and Filtering
func GetProducts(c *fiber.Ctx) error {
	
	// --- (Soal 9) Logika Pagination ---
	
	// Ambil query param 'page' dan 'limit', beri nilai default
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Hitung offset (data yang di-skip)
	offset := (page - 1) * limit

	// --- (Soal 10) Logika Filtering ---
	
	// Buat query dasar
	query := config.DB.Model(&models.Produk{})

	// Filter berdasarkan kategori_id (jika ada)
	if kategoriID := c.Query("kategori_id"); kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}

	// Filter berdasarkan pencarian (search)
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// --- Eksekusi Query ---
	
	var products []models.Produk
	// Terapkan pagination (Offset, Limit) dan filter (dari 'query')
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	
	// (Penting untuk Pagination) Hitung total data
	var total int64
	query.Count(&total) // GORM akan menghitung total baris berdasarkan filter

	// Format response
	return c.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"limit":     limit,
			"last_page": (total + int64(limit) - 1) / int64(limit), // Kalkulasi halaman terakhir
		},
	})
}
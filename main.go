package main

import (
	"ecommerce-api/config"
	"ecommerce-api/handlers"
	"ecommerce-api/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Inisialisasi Database
	config.ConnectDB()

	// Inisiasi Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Static folder untuk file upload
	app.Static("/uploads", "./uploads")

	// Grup Rute /api
	api := app.Group("/api")

	// === Rute Publik (Autentikasi) ===
	api.Post("/register", handlers.Register) // (Soal 1 & 2)
	api.Post("/login", handlers.Login)       // (Soal 1)

	// === Rute Terproteksi (User) ===
	// (Ketentuan 3) Menggunakan JWTMiddleware
	protected := api.Use(middleware.JWTMiddleware())

	// Rute User (Soal 3 & 4)
	userGroup := protected.Group("/user")
	userGroup.Get("/profile", handlers.GetProfile)
	userGroup.Put("/profile", handlers.UpdateProfile)

	// Rute Alamat (Soal 5)
	addressGroup := protected.Group("/addresses")
	addressGroup.Post("/", handlers.CreateAddress)
	addressGroup.Get("/", handlers.GetAddresses)
	addressGroup.Put("/:id", handlers.UpdateAddress)
	addressGroup.Delete("/:id", handlers.DeleteAddress)

	// Rute Produk (Soal 7)
	productGroup := protected.Group("/products")
	productGroup.Post("/", handlers.CreateProduct) // (Ketentuan 5 & 14)
	productGroup.Get("/", handlers.GetProducts) //baru


	// Rute Transaksi (Soal 8)
	protected.Post("/transactions", handlers.Checkout)

	// === Rute Khusus Admin ===
	// (Ketentuan 8) Menggunakan JWTMiddleware + AdminMiddleware
	adminGroup := protected.Group("/admin", middleware.AdminMiddleware())

	// Rute Kategori (Soal 6)
	categoryGroup := adminGroup.Group("/categories")
	categoryGroup.Post("/", handlers.CreateCategory)
	categoryGroup.Get("/", handlers.GetCategories)
	categoryGroup.Get("/:id", handlers.GetCategory)
	categoryGroup.Put("/:id", handlers.UpdateCategory)
	categoryGroup.Delete("/:id", handlers.DeleteCategory)

	log.Fatal(app.Listen(":3000"))
}
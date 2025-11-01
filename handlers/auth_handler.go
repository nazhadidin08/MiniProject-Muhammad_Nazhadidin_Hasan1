package handlers

import (
	"ecommerce-api/config"
	"ecommerce-api/middleware"
	"ecommerce-api/models"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// DTO untuk input register
type RegisterInput struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	NoTelepon string `json:"no_telepon" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// DTO untuk input login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// (Soal 1 & 2) Register User dan otomatis buat Toko
func Register(c *fiber.Ctx) error {
	input := new(RegisterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Validasi input (bisa ditambahkan library validator)

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
		NoTelepon: input.NoTelepon,
	}
	if err := user.HashPassword(input.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// (Soal 8) Gunakan GORM Transaction
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Buat User
		if err := tx.Create(&user).Error; err != nil {
			// (Error 1062: Duplicate entry)
			return err
		}

		// 2. Buat Toko (Soal 2)
		toko := models.Toko{
			UserID: user.ID,
			Name:   user.Name + "'s Shop", // Nama default
		}
		if err := tx.Create(&toko).Error; err != nil {
			return err
		}

		return nil // Commit
	})

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Failed to register user (email may already exist)"})
	}

	// Response DTO
	type UserResponse struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		NoTelepon string `json:"no_telepon"`
	}

	return c.Status(fiber.StatusCreated).JSON(UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		NoTelepon: user.NoTelepon,
	})
}

// (Soal 1) Login User
func Login(c *fiber.Ctx) error {
	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	if err := config.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if err := user.CheckPassword(input.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// (Ketentuan 4) Generate JWT Token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
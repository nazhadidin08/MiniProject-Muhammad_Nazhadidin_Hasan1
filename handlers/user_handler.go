package handlers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"github.com/gofiber/fiber/v2"
)

// Response DTO untuk profil
type UserProfileResponse struct {
	ID    uint          `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
	NoTelepon string   `json:"no_telepon"`
	Role  string        `json:"role"`
	Toko  *models.Toko `json:"toko"` // (Soal 3) Sertakan data toko
}

// (Soal 3) Get User Profile
func GetProfile(c *fiber.Ctx) error {
	// Ambil user_id dari JWT middleware
	userID := c.Locals("user_id").(uint)

	var user models.User
	// (Soal 3) Gunakan GORM Preload untuk menyertakan Toko
	if err := config.DB.Preload("Toko").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	response := UserProfileResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		NoTelepon: user.NoTelepon,
		Role:  user.Role,
		Toko:  user.Toko,
	}

	return c.JSON(response)
}

// DTO untuk update profil
type UpdateProfileInput struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
	NoTelepon string `json:"no_telepon"`
}

// (Soal 4) Update User Profile
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	input := new(UpdateProfileInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Buat map untuk update agar hanya field yang diisi yang di-update
	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["Name"] = input.Name
	}
	if input.Email != "" {
		updates["Email"] = input.Email
	}

	if input.NoTelepon != "" {
		updates["NoTelepon"] = input.NoTelepon
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	// Ambil data terbaru (terutama jika email diubah)
	var updatedUser models.User
	config.DB.First(&updatedUser, userID)
	
	type UpdateResponse struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	return c.JSON(UpdateResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	})
}
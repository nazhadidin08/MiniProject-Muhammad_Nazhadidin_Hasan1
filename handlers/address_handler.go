package handlers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DTO untuk Create/Update Alamat
type AddressInput struct {
	Street     string `json:"street" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}

// (Soal 5 & 12) Create Address
func CreateAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	input := new(AddressInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	address := models.Alamat{
		UserID:     userID, // (Ketentuan 12) Pastikan milik user yg login
		Street:     input.Street,
		City:       input.City,
		Province:   input.Province,
		PostalCode: input.PostalCode,
	}

	if err := config.DB.Create(&address).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create address"})
	}

	return c.Status(fiber.StatusCreated).JSON(address)
}

// (Soal 5 & 12) Get Addresses
func GetAddresses(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var addresses []models.Alamat

	// (Ketentuan 12) Hanya ambil alamat milik user yg login
	config.DB.Where("user_id = ?", userID).Find(&addresses)

	return c.JSON(addresses)
}

// (Soal 5 & 12) Update Address
func UpdateAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	addressID := c.Params("id")

	input := new(AddressInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var address models.Alamat
	// (Ketentuan 12) Pastikan user hanya update alamat miliknya
	if err := config.DB.First(&address, "id = ? AND user_id = ?", addressID, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Address not found or not owned"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	address.Street = input.Street
	address.City = input.City
	address.Province = input.Province
	address.PostalCode = input.PostalCode

	config.DB.Save(&address)
	return c.JSON(address)
}

// (Soal 5 & 12) Delete Address
func DeleteAddress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	addressID := c.Params("id")

	// (Ketentuan 12) Hapus hanya jika ID dan UserID cocok
	result := config.DB.Delete(&models.Alamat{}, "id = ? AND user_id = ?", addressID, userID)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Address not found or not owned"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
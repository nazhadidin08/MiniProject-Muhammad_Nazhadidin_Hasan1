package handlers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"github.com/gofiber/fiber/v2"
)

// (Soal 6) Semua endpoint di sini harus diproteksi AdminMiddleware

// DTO untuk Kategori
type CategoryInput struct {
	Name string `json:"name" validate:"required"`
}

// (Soal 6) Create Category (Admin Only)
func CreateCategory(c *fiber.Ctx) error {
	input := new(CategoryInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	category := models.Kategori{Name: input.Name}
	if err := config.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Category name already exists"})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// (Soal 6) Get All Categories
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Kategori
	config.DB.Find(&categories)
	return c.JSON(categories)
}

// (Soal 6) Get Category By ID
func GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.Kategori
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(category)
}

// (Soal 6) Update Category (Admin Only)
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	input := new(CategoryInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var category models.Kategori
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}

	category.Name = input.Name
	if err := config.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Category name already exists"})
	}

	return c.JSON(category)
}

// (Soal 6) Delete Category (Admin Only)
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.Kategori{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
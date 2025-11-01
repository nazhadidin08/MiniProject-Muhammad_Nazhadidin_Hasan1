package handlers

import (
	"ecommerce-api/config"
	"ecommerce-api/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DTO untuk input Checkout
type CheckoutItemInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
type CheckoutInput struct {
	Items []CheckoutItemInput `json:"items"`
}

// (Soal 8) Checkout (Create Transaction)
func Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	input := new(CheckoutInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if len(input.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cart is empty"})
	}

	var newTransaction models.Transaksi
	var totalPrice float64

	// (Soal 8) Gunakan Transaksi Database GORM (.Transaction())
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Buat header Transaksi (Invoice)
		trx := models.Transaksi{
			UserID:     userID,
			Status:     "pending",
			TotalPrice: 0, // Akan di-update nanti
		}
		if err := tx.Create(&trx).Error; err != nil {
			return err
		}

		totalPrice = 0.0

		for _, item := range input.Items {
			var product models.Produk

			// (Soal 8) Kunci row produk untuk update (mencegah race condition)
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, item.ProductID).Error; err != nil {
				return errors.New("product not found")
			}

			// Cek stok
			if product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for product: %s", product.Name)
			}

			// 2. Kurangi stok di tabel produk
			newStock := product.Stock - item.Quantity
			if err := tx.Model(&product).Update("stock", newStock).Error; err != nil {
				return err
			}

			// 3. (Ketentuan 15 & 16) Buat entri baru di log_produk
			log := models.LogProduk{
				TransaksiID: trx.ID,
				ProdukID:    product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price, // Menyalin harga saat itu
			}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}

			// 4. Buat data detail_trx
			subTotal := product.Price * float64(item.Quantity)
			detail := models.DetailTransaksi{
				TransaksiID: trx.ID,
				ProdukID:    product.ID,
				Quantity:    item.Quantity,
				SubTotal:    subTotal,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			totalPrice += subTotal
		}

		// 5. Update TotalPrice di header Transaksi
		if err := tx.Model(&trx).Update("total_price", totalPrice).Error; err != nil {
			return err
		}

		// Simpan data transaksi yang sudah lengkap untuk response
		// Preload detailnya untuk response
		tx.Preload("DetailTransaksi").Preload("LogProduk").First(&newTransaction, trx.ID)

		return nil // Commit transaksi
	})

	if err != nil {
		// Transaksi otomatis di-rollback jika ada error
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Transaction failed",
			"cause": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newTransaction)
}
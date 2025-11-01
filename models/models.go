package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// (Soal 1)
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	NoTelepon string    `gorm:"uniqueIndex" json:"no_telepon"`
	Password  string    `json:"-"` // Sembunyikan dari JSON output
	Role      string    `gorm:"default:'user'" json:"role"` // 'user' atau 'admin'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi
	Toko      *Toko       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"toko,omitempty"`                 // (Soal 3) User Punya 1 Toko
	Alamat    []Alamat    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"alamat,omitempty"`               // User Punya Banyak Alamat
	Transaksi []Transaksi `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"transaksi,omitempty"`           // User Punya Banyak Transaksi
}

// HashPassword mengenkripsi password sebelum disimpan
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword memvalidasi password
func (u *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
}

// (Soal 2) Toko (Shop)
type Toko struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex" json:"user_id"` // 1-to-1 dengan User
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User   *User    `gorm:"foreignKey:UserID" json:"-"` // Relasi Belongs To
	Produk []Produk `gorm:"foreignKey:TokoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"produk,omitempty"`
}

// (Soal 5) Alamat (Address)
type Alamat struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Street      string    `json:"street"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// (Soal 6) Kategori (Category)
type Kategori struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Produk []Produk `gorm:"foreignKey:KategoriID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"produk,omitempty"`
}

// (Soal 7) Produk (Product)
type Produk struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TokoID      uint      `gorm:"index" json:"toko_id"`     // Relasi ke Toko
	KategoriID  uint      `gorm:"index" json:"kategori_id"` // Relasi ke Kategori
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImagePath   string    `json:"image_path"` // Path ke gambar (Ketentuan 5)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Toko     *Toko     `gorm:"foreignKey:TokoID" json:"-"`
	Kategori *Kategori `gorm:"foreignKey:KategoriID" json:"-"`
}

// (Soal 8) Transaksi (Transaction)
type Transaksi struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"` // e.g., "pending", "paid", "shipped"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User            *User              `gorm:"foreignKey:UserID" json:"-"`
	DetailTransaksi []DetailTransaksi  `gorm:"foreignKey:TransaksiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"detail_transaksi"`
	LogProduk       []LogProduk        `gorm:"foreignKey:TransaksiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"log_produk"`
}

// (Soal 8) DetailTransaksi (TransactionDetail)
type DetailTransaksi struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	TransaksiID uint    `gorm:"index" json:"transaksi_id"`
	ProdukID    uint    `gorm:"index" json:"produk_id"`
	Quantity    int     `json:"quantity"`
	SubTotal    float64 `json:"sub_total"` // PriceAtPurchase * Quantity

	Transaksi *Transaksi `gorm:"foreignKey:TransaksiID" json:"-"`
	Produk    *Produk    `gorm:"foreignKey:ProdukID" json:"-"`
}

// (Soal 15 & 16) LogProduk (ProductLog)
// Menyalin data produk pada saat transaksi terjadi
type LogProduk struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	TransaksiID uint    `gorm:"index" json:"transaksi_id"`
	ProdukID    uint    `gorm:"index" json:"produk_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"` // Harga saat dibeli

	Transaksi *Transaksi `gorm:"foreignKey:TransaksiID" json:"-"`
	Produk    *Produk    `gorm:"foreignKey:ProdukID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
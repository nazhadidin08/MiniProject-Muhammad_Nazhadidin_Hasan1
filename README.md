# E-commerce API Backend (Go-Fiber-GORM)

Ini adalah API backend lengkap untuk platform e-commerce, dibangun menggunakan Go, Fiber, GORM, dan MySQL. Proyek ini mencakup autentikasi JWT, manajemen produk, sistem checkout, dan panel admin.

---

## 🚀 Fitur Utama

* **Autentikasi:** Register, Login, dan proteksi rute menggunakan **JWT (JSON Web Tokens)**.
* **Manajemen User & Toko:** User mendapatkan toko otomatis saat mendaftar.
* **CRUD Alamat:** User dapat mengelola banyak alamat pengiriman.
* **Admin Panel:** *Endpoint* khusus admin untuk mengelola Kategori.
* **Manajemen Produk:** Membuat produk baru dengan **upload file (gambar)**.
* **Sistem Checkout:** Logika checkout kompleks menggunakan **Transaksi Database GORM** untuk memastikan integritas data (mengurangi stok, membuat log produk).
* **Filtering & Pagination:** *Endpoint* untuk mengambil data produk dilengkapi dengan *query parameter* untuk `search`, `kategori_id`, `page`, dan `limit`.
* **Keamanan:** *Middleware* untuk proteksi rute (memerlukan login) dan proteksi *role* (khusus admin).

---

## 🔧 Teknologi yang Digunakan

* **Bahasa:** Go (Golang)
* **Framework:** Fiber (v2)
* **ORM:** GORM
* **Database:** MySQL
* **Autentikasi:** `golang-jwt/jwt`
* **Password Hashing:** `golang.org/x/crypto/bcrypt`
* **Konfigurasi:** `joho/godotenv`

---

## 📂 Struktur Proyek

---

## 📦 Instalasi & Menjalankan

### 1. Prasyarat

* [Go](https://golang.org/dl/) (versi 1.18 atau lebih baru)
* [MySQL](https://dev.mysql.com/downloads/installer/) (atau XAMPP/Docker)
* [Insomnia](https://insomnia.rest/download) atau [Postman](https://www.postman.com/downloads/) (untuk pengujian API)

### 2. Langkah-langkah

1.  **Clone atau Unduh Proyek**
    ```bash
    # (Jika Anda sudah punya reponya)
    git clone [https://github.com/](https://github.com/)[NAMA_ANDA]/ecommerce-api.git
    cd ecommerce-api
    ```

2.  **Buat Database**
    Buka MySQL (phpMyAdmin/HeidiSQL/dll) dan buat database baru dengan nama:
    `db_ecommerce`

3.  **Konfigurasi `.env`**
    Buat file bernama `.env` di *root* proyek dan isi dengan kredensial Anda.

    **.env.template**
    ```.env
    # Sesuaikan dengan user:password dan nama database Anda
    DB_DSN="root:password@tcp(127.0.0.1:3306)/db_ecommerce?charset=utf8mb4&parseTime=True&loc=Local"

    # Ganti dengan KUNCI RAHASIA yang kuat untuk JWT
    JWT_SECRET="kunci_rahasia_anda_yang_sangat_kuat"
    ```

4.  **Instal Dependensi**
    Buka terminal dan jalankan `go mod tidy` untuk mengunduh semua *library* (Fiber, GORM, dll).
    ```bash
    go mod tidy
    ```

5.  **Jalankan Server**
    Server akan berjalan di `http://localhost:3000`.
    ```bash
    go run main.go
    ```
    GORM akan secara otomatis menjalankan `AutoMigrate` dan membuat semua tabel saat server pertama kali dijalankan.

---

## 📚 Dokumentasi API

Semua *endpoint* di bawah `/api` merespons dalam format **JSON**.

### 🔑 Autentikasi

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/register` | Mendaftarkan user baru (dan membuat toko baru). |
| `POST` | `/api/login` | Login user untuk mendapatkan Token JWT. |

### 👤 User (Terproteksi)

*(Memerlukan Bearer Token)*

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/user/profile` | Mengambil profil user yang sedang login (termasuk data Toko). |
| `PUT` | `/api/user/profile` | Memperbarui profil (nama, email, no_telepon) user yang sedang login. |

### 📍 Alamat (Terproteksi)

*(Memerlukan Bearer Token)*

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/addresses` | Membuat alamat pengiriman baru untuk user. |
| `GET` | `/api/addresses` | Melihat semua alamat milik user. |
| `PUT` | `/api/addresses/:id` | Memperbarui alamat milik user. |
| `DELETE` | `/api/addresses/:id` | Menghapus alamat milik user. |

### 📦 Produk (Terproteksi)

*(Memerlukan Bearer Token)*

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/products` | **(Multipart/Form-Data)**. Membuat produk baru & upload gambar. |
| `GET` | `/api/products` | Mengambil semua produk dengan pagination dan filtering. |

**Query Params untuk `GET /api/products`:**
* `page`: Nomor halaman (default: `1`)
* `limit`: Jumlah item per halaman (default: `10`)
* `kategori_id`: Filter berdasarkan ID Kategori.
* `search`: Filter berdasarkan nama/deskripsi produk.

### 💳 Transaksi (Terproteksi)

*(Memerlukan Bearer Token)*

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/transactions` | Melakukan "Checkout". Mengurangi stok, membuat invoice, dan mencatat log. |

### 🛡️ Admin (Khusus Admin)

*(Memerlukan Bearer Token dengan role Admin)*

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/admin/categories` | Membuat Kategori baru. |
| `GET` | `/api/admin/categories` | Mengambil semua Kategori. |
| `GET` | `/api/admin/categories/:id` | Mengambil Kategori berdasarkan ID. |
| `PUT` | `/api/admin/categories/:id` | Memperbarui Kategori. |
| `DELETE` | `/api/admin/categories/:id` | Menghapus Kategori. |
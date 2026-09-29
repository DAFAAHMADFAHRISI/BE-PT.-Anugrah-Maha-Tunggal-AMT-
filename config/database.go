package config

import (
	"fmt"
	"log"
	"os"

	"erp-anugrah-maha-tunggal/internal/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// Muat file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "erp_amt_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v. Pastikan MySQL sudah aktif dan database '%s' sudah dibuat.", err, dbName)
	}

	// Auto-Migrate schema GORM
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Customer{},
		&domain.ForkliftUnit{},
		&domain.Operator{},
		&domain.RentalOrder{},
		&domain.DeliveryLetter{},
		&domain.UnitMovementLog{},
		&domain.Invoice{},
	)
	if err != nil {
		log.Printf("Peringatan saat migrasi tabel GORM: %v", err)
	} else {
		log.Println("Sukses: Migrasi database dan tabel berhasil diselaraskan.")
	}

	// Pastikan kolom departure_time bertipe VARCHAR(50) agar mendukung teks format jam bebas atau kosong
	db.Exec("ALTER TABLE delivery_letters MODIFY COLUMN departure_time VARCHAR(50) NULL")

	// Auto-seed data awal jika belum ada
	seedInitialData(db)

	DB = db
	return db
}

func seedInitialData(db *gorm.DB) {
	defaultHash := "$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6"

	defaultUsers := []domain.User{
		{
			Name:         "Ghofur (Direktur)",
			Username:     "direktur",
			Email:        "Ghofur@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleDirektur,
			Phone:        "081122334455",
			IsActive:     true,
		},
		{
			Name:         "Dafa (Admin Staff)",
			Username:     "admin",
			Email:        "dafa@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleAdminStaff,
			Phone:        "081234567890",
			IsActive:     true,
		},
		{
			Name:         "Luluk (Operasional)",
			Username:     "operasional",
			Email:        "luluk@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleOperasional,
			Phone:        "085611223344",
			IsActive:     true,
		},
		{
			Name:         "Risal (Finance)",
			Username:     "finance",
			Email:        "risal@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleFinance,
			Phone:        "081377889900",
			IsActive:     true,
		},
	}

	for _, u := range defaultUsers {
		var exists int64
		db.Model(&domain.User{}).Where("username = ?", u.Username).Count(&exists)
		if exists == 0 {
			db.Create(&u)
			log.Printf("Seeder: User '%s' (%s) berhasil dibuat.\n", u.Username, u.Role)
		}
	}

	// Seed Unit Forklift jika kosong
	var unitCount int64
	db.Model(&domain.ForkliftUnit{}).Count(&unitCount)
	if unitCount == 0 {
		units := []domain.ForkliftUnit{
			{UnitCode: "FL-01", Brand: "Toyota", Model: "8FD30", CapacityTon: 3.0, FuelType: "Diesel", ManufactureYear: 2021, DailyRate: 1200000, Status: "AVAILABLE", Notes: "Kondisi prima di pool"},
			{UnitCode: "FL-02", Brand: "Komatsu", Model: "FD50AY-10", CapacityTon: 5.0, FuelType: "Diesel", ManufactureYear: 2020, DailyRate: 1800000, Status: "RENTED", Notes: "Sedang bertugas di proyek Cilegon"},
			{UnitCode: "FL-03", Brand: "Mitsubishi", Model: "FD70N", CapacityTon: 7.0, FuelType: "Diesel", ManufactureYear: 2019, DailyRate: 2500000, Status: "AVAILABLE", Notes: "Siap beroperasi"},
		}
		db.Create(&units)
	}

	// Seed Pelanggan jika kosong
	var custCount int64
	db.Model(&domain.Customer{}).Count(&custCount)
	if custCount == 0 {
		customers := []domain.Customer{
			{CustomerCode: "CUST-001", Name: "PT. Pelabuhan Samudera Raya", Type: "PT", PicName: "Pak Rahmat", Phone: "081234567890", Email: "purchasing@samuderaraya.co.id", Address: "Kawasan Industri Cilegon Kav 12"},
			{CustomerCode: "CUST-002", Name: "PT. Mega Baja Mandiri", Type: "PT", PicName: "Ibu Maya", Phone: "081398765432", Email: "logistik@megabaja.com", Address: "Jl. Raya Narogong Km 14, Bekasi"},
		}
		db.Create(&customers)
	}

	// Seed Operator jika kosong
	var oprCount int64
	db.Model(&domain.Operator{}).Count(&oprCount)
	if oprCount == 0 {
		operators := []domain.Operator{
			{Nip: "OPR-01", Name: "Agus Prasetyo", Phone: "081299887766", SioNumber: "SIO-K3-2022-0091", Status: "ON_DUTY"},
			{Nip: "OPR-02", Name: "Doni Kurniawan", Phone: "081344556677", SioNumber: "SIO-K3-2023-0145", Status: "STANDBY"},
		}
		db.Create(&operators)
	}
}

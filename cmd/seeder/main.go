package main

import (
	"flag"
	"fmt"
	"log"

	"erp-anugrah-maha-tunggal/config"
	"erp-anugrah-maha-tunggal/internal/domain"
	"gorm.io/gorm"
)

func main() {
	actionFlag := flag.String("action", "", "Aksi yang diinginkan: 'seed', 'clean', atau 'reset'")
	cleanFlag := flag.Bool("clean", false, "Hapus semua data dalam tabel")
	seedFlag := flag.Bool("seed", false, "Isi data seeder awal")
	resetFlag := flag.Bool("reset", false, "Hapus semua data lalu isi ulang seeder")
	flag.Parse()

	// Menentukan aksi yang dipilih (bisa dari flag maupun argumen langsung seperti: go run cmd/seeder/main.go clean)
	action := "seed" // default

	if len(flag.Args()) > 0 {
		action = flag.Args()[0]
	} else if *cleanFlag {
		action = "clean"
	} else if *resetFlag {
		action = "reset"
	} else if *seedFlag {
		action = "seed"
	} else if *actionFlag != "" {
		action = *actionFlag
	}

	db := config.InitDB()

	switch action {
	case "clean":
		cleanAllData(db)
		fmt.Println("✅ [SUKSES] Semua data dalam database berhasil dibersihkan/dikosongkan.")
	case "seed":
		seedAllData(db)
		fmt.Println("✅ [SUKSES] Data seeder awal berhasil dimasukkan ke dalam database.")
	case "reset":
		cleanAllData(db)
		seedAllData(db)
		fmt.Println("✅ [SUKSES] Database berhasil di-RESET (dikosongkan dan diisi ulang data awal).")
	default:
		log.Fatalf("❌ Aksi '%s' tidak dikenali! Gunakan salah satu: 'clean', 'seed', atau 'reset'\nContoh: go run cmd/seeder/main.go clean", action)
	}
}

// cleanAllData mengosongkan semua tabel dengan aman (mematikan foreign key check sementara)
func cleanAllData(db *gorm.DB) {
	fmt.Println("⏳ Membersihkan seluruh data tabel...")

	tables := []string{
		"unit_movement_logs",
		"invoices",
		"delivery_letters",
		"rental_orders",
		"forklift_units",
		"operators",
		"customers",
		"users",
	}

	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", table)).Error; err != nil {
			log.Printf("Peringatan: Gagal truncate tabel %s: %v", table, err)
		} else {
			fmt.Printf("   - Tabel '%s' berhasil dikosongkan.\n", table)
		}
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
}

// seedAllData mengisi data dummy awal lengkap untuk testing/koreksi
func seedAllData(db *gorm.DB) {
	fmt.Println("⏳ Memasukkan data seeder awal...")

	// 1. Seed Users (Bcrypt hash untuk 'password123')
	defaultHash := "$2a$10$Ndq4ce69r0nlmsG25fLfyOAj9muridcFWovko10gFw11HV08H5Tci"
	users := []domain.User{
		// Direktur
		{
			Name:         "Bapak Ir. Hendra Prasetyo (Direktur Utama)",
			Username:     "direktur",
			Email:        "direktur@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleDirektur,
			Phone:        "081122334455",
			IsActive:     true,
		},
		{
			Name:         "Bapak Bambang Wijaya (Direktur Operasional)",
			Username:     "direktur_ops",
			Email:        "bambang.direksi@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleDirektur,
			Phone:        "081188776655",
			IsActive:     true,
		},
		// Admin Staff
		{
			Name:         "Siti Rahma, S.Kom (Admin Staff Rental)",
			Username:     "admin",
			Email:        "admin@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleAdminStaff,
			Phone:        "081234567890",
			IsActive:     true,
		},
		{
			Name:         "Anisa Putri, A.Md (Admin Kontrak & PO)",
			Username:     "admin2",
			Email:        "anisa.admin@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleAdminStaff,
			Phone:        "081298765432",
			IsActive:     true,
		},
		// Operasional
		{
			Name:         "Budi Santoso (Koordinator Lapangan & Armada)",
			Username:     "operasional",
			Email:        "operasional@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleOperasional,
			Phone:        "085611223344",
			IsActive:     true,
		},
		{
			Name:         "Agus Setiawan (Dispatcher & Pengawas Lapangan)",
			Username:     "operasional2",
			Email:        "agus.ops@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleOperasional,
			Phone:        "085722334455",
			IsActive:     true,
		},
		// Finance
		{
			Name:         "Dewi Lestari, S.E. (Head of Finance & Billing)",
			Username:     "finance",
			Email:        "finance@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleFinance,
			Phone:        "081377889900",
			IsActive:     true,
		},
		{
			Name:         "Rina Marlina, S.Ak (Staff Finance & Penagihan)",
			Username:     "finance2",
			Email:        "rina.finance@anugrahmahatunggal.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleFinance,
			Phone:        "081366554433",
			IsActive:     true,
		},
	}

	for _, u := range users {
		var exists int64
		db.Model(&domain.User{}).Where("username = ?", u.Username).Count(&exists)
		if exists == 0 {
			db.Create(&u)
		}
	}
	fmt.Printf("   - %d Akun User berhasil disiapkan.\n", len(users))

	// 2. Seed Customer
	customers := []domain.Customer{
		{CustomerCode: "CUST-001", Name: "PT. Pelabuhan Samudera Raya", Type: "PT", PicName: "Pak Rahmat", Phone: "081234567890", Email: "purchasing@samuderaraya.co.id", Address: "Kawasan Industri Cilegon Kav 12, Banten"},
		{CustomerCode: "CUST-002", Name: "PT. Mega Baja Mandiri", Type: "PT", PicName: "Ibu Maya", Phone: "081398765432", Email: "logistik@megabaja.com", Address: "Jl. Raya Narogong Km 14, Bekasi"},
		{CustomerCode: "CUST-003", Name: "CV. Makmur Jaya Abadi", Type: "CV", PicName: "Pak Joko", Phone: "085711223344", Email: "makmurjaya@gmail.com", Address: "Jl. Industri Pergudangan No. 45, Tangerang"},
	}
	for _, c := range customers {
		var exists int64
		db.Model(&domain.Customer{}).Where("customer_code = ?", c.CustomerCode).Count(&exists)
		if exists == 0 {
			db.Create(&c)
		}
	}
	fmt.Printf("   - %d Pelanggan (Customer) berhasil disiapkan.\n", len(customers))

	// 3. Seed Forklift Units
	units := []domain.ForkliftUnit{
		{UnitCode: "FL-01", Brand: "Toyota", Model: "8FD30", CapacityTon: 3.0, FuelType: "Diesel", ManufactureYear: 2021, DailyRate: 1200000, Status: "AVAILABLE", Notes: "Kondisi prima, siap operasi"},
		{UnitCode: "FL-02", Brand: "Komatsu", Model: "FD50AY-10", CapacityTon: 5.0, FuelType: "Diesel", ManufactureYear: 2020, DailyRate: 1800000, Status: "RENTED", Notes: "Sedang bertugas di proyek Samudera Raya"},
		{UnitCode: "FL-03", Brand: "Mitsubishi", Model: "FD70N", CapacityTon: 7.0, FuelType: "Diesel", ManufactureYear: 2019, DailyRate: 2500000, Status: "AVAILABLE", Notes: "Kapasitas 7 ton, standby di pool"},
		{UnitCode: "FL-04", Brand: "Toyota", Model: "8FB25", CapacityTon: 2.5, FuelType: "Electric", ManufactureYear: 2022, DailyRate: 1100000, Status: "AVAILABLE", Notes: "Unit baterai ramah lingkungan untuk gudang indoor"},
		{UnitCode: "FL-05", Brand: "TCM", Model: "FD100-2", CapacityTon: 10.0, FuelType: "Diesel", ManufactureYear: 2018, DailyRate: 3800000, Status: "MAINTENANCE", Notes: "Jadwal servis hidrolik rutin"},
	}
	for _, u := range units {
		var exists int64
		db.Model(&domain.ForkliftUnit{}).Where("unit_code = ?", u.UnitCode).Count(&exists)
		if exists == 0 {
			db.Create(&u)
		}
	}
	fmt.Printf("   - %d Unit Forklift berhasil disiapkan.\n", len(units))

	// 4. Seed Operators
	operators := []domain.Operator{
		{Nip: "OPR-01", Name: "Agus Prasetyo", Phone: "081299887766", SioNumber: "SIO-K3-2022-0091", Status: "ON_DUTY"},
		{Nip: "OPR-02", Name: "Doni Kurniawan", Phone: "081344556677", SioNumber: "SIO-K3-2023-0145", Status: "STANDBY"},
		{Nip: "OPR-03", Name: "Rian Hidayat", Phone: "085677889900", SioNumber: "SIO-K3-2021-0032", Status: "STANDBY"},
	}
	for _, o := range operators {
		var exists int64
		db.Model(&domain.Operator{}).Where("nip = ?", o.Nip).Count(&exists)
		if exists == 0 {
			db.Create(&o)
		}
	}
	fmt.Printf("   - %d Operator Forklift berhasil disiapkan.\n", len(operators))

	// 5. Seed Contoh Transaksi Pesanan & Surat Jalan Aktif
	var orderCount int64
	db.Model(&domain.RentalOrder{}).Where("order_number = ?", "ORD-202609-0001").Count(&orderCount)
	if orderCount == 0 {
		var cust domain.Customer
		db.Where("customer_code = ?", "CUST-001").First(&cust)

		var flUnit domain.ForkliftUnit
		db.Where("unit_code = ?", "FL-02").First(&flUnit)

		var adminUser domain.User
		db.Where("username = ?", "admin").First(&adminUser)

		unitID := flUnit.ID
		createdBy := adminUser.ID
		endDate := "2026-09-26"

		sampleOrder := domain.RentalOrder{
			OrderNumber:         "ORD-202609-0001",
			CustomerID:          cust.ID,
			OrderDate:           "2026-09-24",
			StartDate:           "2026-09-24",
			EndDate:             &endDate,
			RentalDurationType:  "HARI",
			DurationValue:       3,
			ProjectLocation:     "Kawasan Industri Cilegon Kav 12, Area Dermaga 3",
			LocationPicName:     "Pak Herman",
			LocationPicPhone:    "081288990011",
			RequiredCapacityTon: 5.0,
			UnitID:              &unitID,
			AgreedPrice:         5400000,
			Status:              "ACTIVE",
			PoReference:         "PO-PSR/IX/2026/88",
			Notes:               "Bongkar kontainer baja plat dermaga",
			CreatedBy:           &createdBy,
		}
		db.Create(&sampleOrder)

		var opr domain.Operator
		db.Where("nip = ?", "OPR-01").First(&opr)

		sampleSJ := domain.DeliveryLetter{
			LetterNumber:      "SJ-202609-0001",
			RentalOrderID:     sampleOrder.ID,
			UnitID:            flUnit.ID,
			OperatorID:        opr.ID,
			IssueDate:         "2026-09-24",
			DepartureTime:     "08:30:00",
			JobDescription:    "Bongkar muat kontainer baja plat dan pemindahan material dermaga",
			RecipientName:     "Pak Herman",
			OperationalStatus: "WORKING",
			Notes:             "Unit dan operator sudah tiba di dermaga",
		}
		db.Create(&sampleSJ)

		fmt.Println("   - Contoh Pesanan Sewa & Surat Jalan aktif berhasil dibuat.")
	}
}

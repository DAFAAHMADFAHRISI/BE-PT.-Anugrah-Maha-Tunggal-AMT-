package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"erp-anugrah-maha-tunggal/config"
	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/service"
	"time"
	"gorm.io/gorm"
)

func main() {
	actionFlag := flag.String("action", "", "Aksi yang diinginkan: 'seed', 'clean', atau 'reset'")
	cleanFlag := flag.Bool("clean", false, "Hapus semua data dalam tabel")
	seedFlag := flag.Bool("seed", false, "Isi data seeder awal")
	resetFlag := flag.Bool("reset", false, "Hapus semua data lalu isi ulang seeder")
	flag.Parse()

	action := "reset" // default reset agar data excel tersinkronisasi bersih

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
		fmt.Println("✅ [SUKSES] Database berhasil di-RESET menggunakan data Excel yang diselaraskan!")
	default:
		log.Fatalf("❌ Aksi '%s' tidak dikenali! Gunakan salah satu: 'clean', 'seed', atau 'reset'", action)
	}
}

// cleanAllData mengosongkan semua tabel dengan aman
func cleanAllData(db *gorm.DB) {
	fmt.Println("⏳ Membersihkan seluruh data tabel...")

	tables := []string{
		"detail_ranking",
		"hasil_ranking",
		"spk_maintenance",
		"log_operasional",
		"bobot_kriteria",
		"kriteria",
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

// seedAllData mengisi data seeder lengkap berbasis dataset Excel PT. Anugrah Maha Tunggal
func seedAllData(db *gorm.DB) {
	fmt.Println("⏳ Memasukkan data seeder berbasis Excel...")

	// 1. Seed Users (Bcrypt hash untuk 'pw123')
	defaultHash := "$2a$10$K2cWUtAqm19gYRsHC4s9dO4m4Tq4xhvRTUYg4yap.hpOmUUV6HBX6"
	users := []domain.User{
		{
			Name:         "Ghofur (Direktur Utama)",
			Username:     "direktur",
			Email:        "Ghofur@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleDirektur,
			Phone:        "081122334455",
			IsActive:     true,
		},
		{
			Name:         "Dafa (Staff Admin Bengkel & Rental)",
			Username:     "admin",
			Email:        "dafa@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleAdminStaff,
			Phone:        "081234567890",
			IsActive:     true,
		},
		{
			Name:         "Luluk (Koordinator Lapangan & Kepala Mekanik)",
			Username:     "operasional",
			Email:        "luluk@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleOperasional,
			Phone:        "085611223344",
			IsActive:     true,
		},
		{
			Name:         "Risal (Staff Keuangan & Billing)",
			Username:     "finance",
			Email:        "risal@gmail.com",
			PasswordHash: defaultHash,
			Role:         domain.RoleFinance,
			Phone:        "081377889900",
			IsActive:     true,
		},
	}

	for i := range users {
		db.Create(&users[i])
	}
	fmt.Printf("   - %d Akun Pengguna berhasil disiapkan.\n", len(users))

	// 2. Seed Customer dari Excel
	customers := []domain.Customer{
		{CustomerCode: "CUST-001", Name: "PT. Metalindo Guna", Type: "PT", PicName: "Pak Rudi", Phone: "081234567891", Email: "metalindo@gmail.com", Address: "Kawasan PIK, Jakarta Utara"},
		{CustomerCode: "CUST-002", Name: "PT. Fajar Bahari", Type: "PT", PicName: "Ibu Siska", Phone: "081234567892", Email: "fajarbahari@gmail.com", Address: "Marunda Center Blok A, Jakarta Utara"},
		{CustomerCode: "CUST-003", Name: "PT. Contena Pratama", Type: "PT", PicName: "Pak Hendra", Phone: "081234567893", Email: "contena@gmail.com", Address: "Sungai Tiram, Cilincing, Jakarta Utara"},
		{CustomerCode: "CUST-004", Name: "PT. Basst", Type: "PT", PicName: "Pak Anton", Phone: "081234567894", Email: "basst@gmail.com", Address: "Kawasan Industri Pulo Gadung, Jakarta Timur"},
		{CustomerCode: "CUST-005", Name: "PT. Tegar Jaya Mandiri", Type: "PT", PicName: "Pak Budi", Phone: "081234567895", Email: "tegarjaya@gmail.com", Address: "Marunda Center Blok C, Jakarta Utara"},
		{CustomerCode: "CUST-006", Name: "PT. Indokomas", Type: "PT", PicName: "Pak Denny", Phone: "081234567896", Email: "indokomas@gmail.com", Address: "Pulo Gadung, Jakarta Timur"},
		{CustomerCode: "CUST-007", Name: "PT. Maju Alam Sutra", Type: "PT", PicName: "Ibu Maya", Phone: "081234567897", Email: "majualam@gmail.com", Address: "Central Cakung, Jakarta Timur"},
		{CustomerCode: "CUST-008", Name: "PT. Global Mas Logistik", Type: "PT", PicName: "Pak Ferry", Phone: "081234567898", Email: "globalmas@gmail.com", Address: "Senayan, Jakarta Pusat"},
		{CustomerCode: "CUST-009", Name: "H. Zaini", Type: "Perorangan", PicName: "H. Zaini", Phone: "081234567899", Email: "zaini@gmail.com", Address: "Pegangsaan II, Kelapa Gading"},
		{CustomerCode: "CUST-010", Name: "Bp. Hariri", Type: "Perorangan", PicName: "Bp. Hariri", Phone: "081234567800", Email: "hariri@gmail.com", Address: "Ramayana Semper, Jakarta Utara"},
		{CustomerCode: "CUST-011", Name: "H. Sai", Type: "Perorangan", PicName: "H. Sai", Phone: "081234567801", Email: "hsai@gmail.com", Address: "Jl. Perintis Kemerdekaan, Jakarta Timur"},
		{CustomerCode: "CUST-012", Name: "Bp. Wahid", Type: "Perorangan", PicName: "Bp. Wahid", Phone: "081234567802", Email: "wahid@gmail.com", Address: "Pegangsaan II, Jakarta Utara"},
	}
	for i := range customers {
		db.Create(&customers[i])
	}
	fmt.Printf("   - %d Pelanggan (Customer) dari data Excel berhasil disiapkan.\n", len(customers))

	// 3. Seed Master Operator dari Excel
	operators := []domain.Operator{
		{Nip: "OPR-01", Name: "Ely", Phone: "081299887711", SioNumber: "SIO-K3-2022-001", Status: "STANDBY"},
		{Nip: "OPR-02", Name: "Roni", Phone: "081299887722", SioNumber: "SIO-K3-2022-002", Status: "ON_DUTY"},
		{Nip: "OPR-03", Name: "Oni", Phone: "081299887733", SioNumber: "SIO-K3-2022-003", Status: "STANDBY"},
		{Nip: "OPR-04", Name: "Deni", Phone: "081299887744", SioNumber: "SIO-K3-2022-004", Status: "STANDBY"},
		{Nip: "OPR-05", Name: "Roji", Phone: "081299887755", SioNumber: "SIO-K3-2022-005", Status: "ON_DUTY"},
		{Nip: "OPR-06", Name: "Umri", Phone: "081299887766", SioNumber: "SIO-K3-2022-006", Status: "STANDBY"},
		{Nip: "OPR-07", Name: "Jasim", Phone: "081299887777", SioNumber: "SIO-K3-2022-007", Status: "STANDBY"},
		{Nip: "OPR-08", Name: "Aris", Phone: "081299887788", SioNumber: "SIO-K3-2022-008", Status: "STANDBY"},
	}
	for i := range operators {
		db.Create(&operators[i])
	}
	fmt.Printf("   - %d Operator Forklift berlisensi dari data Excel berhasil disiapkan.\n", len(operators))

	// 4. Seed 10 Armada Forklift dari Excel lengkap dengan variabel mentah C1..C4
	// C1: Total Jam Kerja, C2: Usia Mesin (Tahun), C3: Frekuensi Rusak (Kali), C4: Sisa Hari Bebas (Hari)
	units := []domain.ForkliftUnit{
		{
			UnitCode:           "FK-2.5-01",
			Brand:              "Toyota",
			Model:              "8FB25",
			CapacityTon:        2.5,
			FuelType:           "Electric",
			ManufactureYear:    2021,
			DailyRate:          700000,
			Status:             "AVAILABLE",
			Notes:              "Unit baterai elektrik ramah lingkungan untuk gudang tertutup",
			HourMeter:          127.0, // C1
			MachineAge:         7.0,   // C2
			BreakdownFrequency: 4.0,   // C3
			IdleDays:           4.0,   // C4
		},
		{
			UnitCode:           "FK-2.5-02",
			Brand:              "Toyota",
			Model:              "8FD25",
			CapacityTon:        2.5,
			FuelType:           "Diesel",
			ManufactureYear:    2020,
			DailyRate:          700000,
			Status:             "AVAILABLE",
			Notes:              "Unit diesel standar operasional lincah di pool",
			HourMeter:          200.0, // C1
			MachineAge:         3.0,   // C2
			BreakdownFrequency: 2.0,   // C3
			IdleDays:           2.0,   // C4
		},
		{
			UnitCode:           "FK-2.5-03",
			Brand:              "Komatsu",
			Model:              "FD25T-17",
			CapacityTon:        2.5,
			FuelType:           "Diesel",
			ManufactureYear:    2022,
			DailyRate:          700000,
			Status:             "AVAILABLE",
			Notes:              "Mesin diesel bertenaga, jam kerja tinggi",
			HourMeter:          194.0, // C1
			MachineAge:         4.0,   // C2
			BreakdownFrequency: 4.0,   // C3
			IdleDays:           7.0,   // C4
		},
		{
			UnitCode:           "FK-03-01",
			Brand:              "Toyota",
			Model:              "8FD30",
			CapacityTon:        3.0,
			FuelType:           "Diesel",
			ManufactureYear:    2019,
			DailyRate:          900000,
			Status:             "AVAILABLE",
			Notes:              "Unit serbaguna kapasitas 3 ton, siap operasi",
			HourMeter:          185.0, // C1
			MachineAge:         9.0,   // C2
			BreakdownFrequency: 3.0,   // C3
			IdleDays:           8.0,   // C4
		},
		{
			UnitCode:           "FK-3.5-01",
			Brand:              "Mitsubishi",
			Model:              "FD35N",
			CapacityTon:        3.5,
			FuelType:           "Diesel",
			ManufactureYear:    2021,
			DailyRate:          1000000,
			Status:             "AVAILABLE",
			Notes:              "Kapasitas 3.5 ton, sering disewa di depo kontainer",
			HourMeter:          196.0, // C1
			MachineAge:         9.0,   // C2
			BreakdownFrequency: 3.0,   // C3
			IdleDays:           7.0,   // C4
		},
		{
			UnitCode:           "FK-04-01",
			Brand:              "Komatsu",
			Model:              "FD40T-10",
			CapacityTon:        4.0,
			FuelType:           "Diesel",
			ManufactureYear:    2018,
			DailyRate:          1200000,
			Status:             "AVAILABLE",
			Notes:              "Kapasitas 4 ton, tangguh untuk industri manufaktur",
			HourMeter:          132.0, // C1
			MachineAge:         8.0,   // C2
			BreakdownFrequency: 4.0,   // C3
			IdleDays:           9.0,   // C4
		},
		{
			UnitCode:           "FK-05-01",
			Brand:              "TCM",
			Model:              "FD50T",
			CapacityTon:        5.0,
			FuelType:           "Diesel",
			ManufactureYear:    2017,
			DailyRate:          1400000,
			Status:             "AVAILABLE",
			Notes:              "Unit kelas 5 ton untuk bongkar muat kargo pelabuhan",
			HourMeter:          171.0, // C1
			MachineAge:         8.0,   // C2
			BreakdownFrequency: 4.0,   // C3
			IdleDays:           6.0,   // C4
		},
		{
			UnitCode:           "FK-07-01",
			Brand:              "Mitsubishi",
			Model:              "FD70N",
			CapacityTon:        7.0,
			FuelType:           "Diesel",
			ManufactureYear:    2019,
			DailyRate:          2500000,
			Status:             "RENTED",
			Notes:              "Kapasitas 7 ton, sedang bertugas di proyek dermaga",
			HourMeter:          149.0, // C1
			MachineAge:         8.0,   // C2
			BreakdownFrequency: 2.0,   // C3
			IdleDays:           10.0,  // C4
		},
		{
			UnitCode:           "FK-07-02",
			Brand:              "Komatsu",
			Model:              "FD70T-10",
			CapacityTon:        7.0,
			FuelType:           "Diesel",
			ManufactureYear:    2020,
			DailyRate:          2500000,
			Status:             "AVAILABLE",
			Notes:              "Unit 7 ton kedua, performa hidrolik prima",
			HourMeter:          106.0, // C1
			MachineAge:         3.0,   // C2
			BreakdownFrequency: 4.0,   // C3
			IdleDays:           1.0,   // C4
		},
		{
			UnitCode:           "FK-10-01",
			Brand:              "TCM",
			Model:              "FD100-2",
			CapacityTon:        10.0,
			FuelType:           "Diesel",
			ManufactureYear:    2016,
			DailyRate:          3500000,
			Status:             "MAINTENANCE",
			Notes:              "Unit heavy duty 10 ton, usia mesin tua dan jam kerja tinggi",
			HourMeter:          174.0, // C1
			MachineAge:         10.0,  // C2
			BreakdownFrequency: 4.0,   // C3
			IdleDays:           9.0,   // C4
		},
	}

	for i := range units {
		db.Create(&units[i])
	}
	fmt.Printf("   - %d Unit Armada Forklift dari data Excel berhasil didaftarkan.\n", len(units))

	// 5. Seed Master 4 Kriteria & Bobot Preferensi SAW
	kriteriaList := []domain.Kriteria{
		{
			Code:          "C1",
			Name:          "Total Jam Kerja (Hour Meter)",
			AttributeType: "benefit",
			Unit:          "Jam",
			Description:   "Akumulasi jam operasional mesin forklift. Semakin tinggi jam kerja, semakin mendesak unit membutuhkan servis/penggantian komponen.",
		},
		{
			Code:          "C2",
			Name:          "Usia Mesin",
			AttributeType: "benefit",
			Unit:          "Tahun",
			Description:   "Usia operasional mesin unit forklift sejak tahun pembuatan. Semakin tua mesin, semakin tinggi potensi aus komponen mesin.",
		},
		{
			Code:          "C3",
			Name:          "Frekuensi Kerusakan",
			AttributeType: "benefit",
			Unit:          "Kali",
			Description:   "Jumlah catatan kendala teknis atau perbaikan dalam 3 bulan terakhir. Semakin sering rusak, semakin tinggi prioritas pemeliharaan.",
		},
		{
			Code:          "C4",
			Name:          "Sisa Hari Bebas Tugas (Idle Time)",
			AttributeType: "benefit",
			Unit:          "Hari",
			Description:   "Jumlah hari unit berada di garasi/pool tanpa kontrak sewa aktif. Semakin banyak hari bebas, semakin optimal kesempatan ditarik ke bengkel.",
		},
	}

	weights := map[string]float64{
		"C1": 0.30, // 30%
		"C2": 0.20, // 20%
		"C3": 0.35, // 35%
		"C4": 0.15, // 15%
	}

	for _, k := range kriteriaList {
		db.Create(&k)
		w := weights[k.Code]
		bobot := domain.BobotKriteria{
			KriteriaID: k.ID,
			Weight:     w,
			Percentage: w * 100,
			Notes:      fmt.Sprintf("Bobot preferensi manajemen untuk kriteria %s sebesar %.0f%%", k.Code, w*100),
		}
		db.Create(&bobot)
	}
	fmt.Println("   - Master 4 Kriteria (C1..C4) dan Bobot Preferensi (W1..W4) berhasil disimpan.")

	// 6. Seed Riwayat Log Operasional Harian awal untuk setiap unit
	staffID := users[1].ID // Dafa (Admin)
	for _, u := range units {
		logOp := domain.OperationalLog{
			UnitID:        u.ID,
			LogDate:       "2026-11-10",
			HourMeter:     u.HourMeter,
			MachineAge:    u.MachineAge,
			BreakdownFreq: u.BreakdownFrequency,
			IdleDays:      u.IdleDays,
			RecordedBy:    &staffID,
			Notes:         fmt.Sprintf("Pencatatan variabel operasional rutin unit %s (%s)", u.UnitCode, u.Brand),
		}
		db.Create(&logOp)
	}
	fmt.Printf("   - %d Log Operasional Armada awal berhasil direkam.\n", len(units))

	// 7. Eksekusi Engine SAW Pertama Kali untuk Menghasilkan Snapshot Ranking Awal
	sawSvc := service.NewSAWService(db)
	calcResult, err := sawSvc.CalculateAndSave(context.Background(), &users[0].ID, "Kalkulasi inisialisasi sistem SPK berbasis dataset Excel")
	if err != nil {
		log.Printf("Peringatan kalkulasi SAW awal: %v", err)
	} else {
		fmt.Printf("   - Engine SAW berhasil menghitung %d unit (Kritis/Merah: %d, Waspada/Kuning: %d, Aman/Hijau: %d).\n",
			calcResult.TotalUnits, calcResult.CriticalCount, calcResult.WarningCount, calcResult.SafeCount)
	}

	// 8. Seed Contoh Transaksi Order Sewa & Surat Jalan Aktif
	orderDate := "2026-11-08"
	startDate := "2026-11-09"
	endDate := "2026-11-12"
	adminID := users[1].ID
	unit7ID := units[7].ID // FK-07-01 (Sedang disewa)

	sampleOrder := domain.RentalOrder{
		OrderNumber:         "ORD-202611-0001",
		CustomerID:          customers[0].ID, // PT. Metalindo Guna
		OrderDate:           orderDate,
		StartDate:           startDate,
		EndDate:             &endDate,
		RentalDurationType:  "HARI",
		DurationValue:       3,
		ProjectLocation:     "Kawasan PIK, Area Pergudangan Logistik Blok B-12",
		LocationPicName:     "Pak Rudi",
		LocationPicPhone:    "081234567891",
		RequiredCapacityTon: 7.0,
		UnitID:              &unit7ID,
		AgreedPrice:         7500000,
		Status:              "ACTIVE",
		PoReference:         "PO-MTL/XI/2026/045",
		Notes:               "Bongkar muatan alat berat dan material baja di dermaga PIK",
		CreatedBy:           &adminID,
	}
	db.Create(&sampleOrder)

	sampleSJ := domain.DeliveryLetter{
		LetterNumber:      "SJ-202611-0001",
		RentalOrderID:     sampleOrder.ID,
		UnitID:            unit7ID,
		OperatorID:        operators[1].ID, // Roni
		IssueDate:         "2026-11-09",
		DepartureTime:     "08:00:00",
		JobDescription:    "Bongkar muat kontainer baja dan pemindahan material kargo",
		RecipientName:     "Pak Rudi",
		OperationalStatus: "WORKING",
		Notes:             "Unit dan operator bertugas di lokasi proyek PIK",
	}
	db.Create(&sampleSJ)

	// 9. Seed Contoh Dokumen SPK Pemeliharaan Bengkel Terbitan SAW
	// Unit FK-10-01 (Kapasitas 10 Ton, Kritis) diterbitkan SPK
	unit10ID := units[9].ID // FK-10-01
	direkturID := users[0].ID
	sampleSPK := domain.SPKMaintenance{
		SPKNumber:      "SPK-MT-20261110-001",
		UnitID:         unit10ID,
		RiskScore:      0.8250,
		RiskStatus:     "CRITICAL",
		MechanicName:   "Bp. Slamet Riyadi (Kepala Mekanik Bengkel)",
		ComplaintNotes: "Overhaul sistem hidrolik silinder utama, kuras oli transmisi, dan penggantian filter solar & oli berkala",
		ScheduledDate:  "2026-11-12",
		Status:         "IN_PROGRESS",
		IssuedBy:       &direkturID,
		IssuedAt:       time.Now(),
		Notes:          "Diterbitkan otomatis via SPK Rekomendasi SAW. Unit berstatus Under Maintenance.",
	}
	db.Create(&sampleSPK)
	fmt.Println("   - Contoh Surat Perintah Kerja (SPK) Bengkel berhasil diterbitkan.")
}

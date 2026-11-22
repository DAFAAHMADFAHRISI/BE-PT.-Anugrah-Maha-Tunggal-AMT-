package service_test

import (
	"testing"

	"erp-anugrah-maha-tunggal/internal/service"
)

// TestComputeSAW_ExcelData menguji algoritma SAW menggunakan data contoh dari Excel transaksi
func TestComputeSAW_ExcelData(t *testing.T) {
	// Kriteria dari sistem SPK
	criteria := []service.SAWCriteria{
		{Code: "C1", Name: "Total Jam Kerja", AttributeType: "benefit", Weight: 0.30, Unit: "Jam"},
		{Code: "C2", Name: "Usia Mesin", AttributeType: "benefit", Weight: 0.20, Unit: "Tahun"},
		{Code: "C3", Name: "Frekuensi Kerusakan", AttributeType: "benefit", Weight: 0.35, Unit: "Kali"},
		{Code: "C4", Name: "Sisa Hari Bebas Tugas", AttributeType: "benefit", Weight: 0.15, Unit: "Hari"},
	}

	// Data contoh alternatif 10 unit forklift dari data Excel
	alternatives := []service.SAWAlternative{
		{
			UnitID:      1,
			UnitCode:    "FK-03-01",
			Brand:       "Toyota",
			Model:       "8FD30",
			CapacityTon: 3.0,
			RawValues: map[string]float64{
				"C1": 62, "C2": 6, "C3": 0, "C4": 3,
			},
		},
		{
			UnitID:      2,
			UnitCode:    "FK-07-01",
			Brand:       "Mitsubishi",
			Model:       "FD70N",
			CapacityTon: 7.0,
			RawValues: map[string]float64{
				"C1": 149, "C2": 8, "C3": 2, "C4": 10,
			},
		},
		{
			UnitID:      3,
			UnitCode:    "FK-07-02",
			Brand:       "Komatsu",
			Model:       "FD70T",
			CapacityTon: 7.0,
			RawValues: map[string]float64{
				"C1": 106, "C2": 3, "C3": 4, "C4": 1,
			},
		},
		{
			UnitID:      4,
			UnitCode:    "FK-2.5-01",
			Brand:       "Toyota",
			Model:       "8FB25",
			CapacityTon: 2.5,
			RawValues: map[string]float64{
				"C1": 98, "C2": 6, "C3": 4, "C4": 7,
			},
		},
		{
			UnitID:      5,
			UnitCode:    "FK-2.5-02",
			Brand:       "Toyota",
			Model:       "8FD25",
			CapacityTon: 2.5,
			RawValues: map[string]float64{
				"C1": 43, "C2": 2, "C3": 2, "C4": 6,
			},
		},
		{
			UnitID:      6,
			UnitCode:    "FK-2.5-03",
			Brand:       "Komatsu",
			Model:       "FD25T",
			CapacityTon: 2.5,
			RawValues: map[string]float64{
				"C1": 194, "C2": 2, "C3": 4, "C4": 7,
			},
		},
		{
			UnitID:      7,
			UnitCode:    "FK-10-01",
			Brand:       "TCM",
			Model:       "FD100",
			CapacityTon: 10.0,
			RawValues: map[string]float64{
				"C1": 123, "C2": 10, "C3": 2, "C4": 9,
			},
		},
	}

	result, err := service.ComputeSAW(alternatives, criteria)
	if err != nil {
		t.Fatalf("ComputeSAW error: %v", err)
	}

	if len(result.Alternatives) != len(alternatives) {
		t.Errorf("Expected %d results, got %d", len(alternatives), len(result.Alternatives))
	}

	// Pastikan urutan peringkat tersortir secara descending (Vi rank 1 >= Vi rank 2)
	for i := 0; i < len(result.Alternatives)-1; i++ {
		current := result.Alternatives[i]
		next := result.Alternatives[i+1]
		if current.ScoreV < next.ScoreV {
			t.Errorf("Ranking sort violation: Rank %d (%.4f) < Rank %d (%.4f)",
				current.RankOrder, current.ScoreV, next.RankOrder, next.ScoreV)
		}
		if current.RankOrder != i+1 {
			t.Errorf("Expected RankOrder %d, got %d", i+1, current.RankOrder)
		}
	}

	// Verifikasi nilai max yang dihitung
	expectedMaxC1 := 194.0
	expectedMaxC2 := 10.0
	expectedMaxC3 := 4.0
	expectedMaxC4 := 10.0

	if result.MaxValues["C1"] != expectedMaxC1 {
		t.Errorf("Expected Max C1 = %.2f, got %.2f", expectedMaxC1, result.MaxValues["C1"])
	}
	if result.MaxValues["C2"] != expectedMaxC2 {
		t.Errorf("Expected Max C2 = %.2f, got %.2f", expectedMaxC2, result.MaxValues["C2"])
	}
	if result.MaxValues["C3"] != expectedMaxC3 {
		t.Errorf("Expected Max C3 = %.2f, got %.2f", expectedMaxC3, result.MaxValues["C3"])
	}
	if result.MaxValues["C4"] != expectedMaxC4 {
		t.Errorf("Expected Max C4 = %.2f, got %.2f", expectedMaxC4, result.MaxValues["C4"])
	}

	// Cek unit peringkat 1 memiliki status risiko dan skor yang valid (0 < Vi <= 1)
	topUnit := result.Alternatives[0]
	if topUnit.ScoreV <= 0 || topUnit.ScoreV > 1.0 {
		t.Errorf("Invalid ScoreV for top unit: %.4f", topUnit.ScoreV)
	}
	if topUnit.RiskStatus == "" || topUnit.RiskColor == "" {
		t.Errorf("RiskStatus or RiskColor is empty")
	}

	t.Logf("Sukses: Peringkat 1 adalah Unit %s (%s) dengan Skor Vi = %.4f [Status: %s - %s]",
		topUnit.UnitCode, topUnit.Brand, topUnit.ScoreV, topUnit.RiskStatus, topUnit.RiskColor)
}

// TestComputeSAW_WeightValidation menguji validasi total bobot != 100%
func TestComputeSAW_WeightValidation(t *testing.T) {
	criteriaInvalid := []service.SAWCriteria{
		{Code: "C1", Weight: 0.30},
		{Code: "C2", Weight: 0.20},
		// Total hanya 0.50 (50%), harus gagal
	}
	alternatives := []service.SAWAlternative{
		{UnitID: 1, UnitCode: "FK-01", RawValues: map[string]float64{"C1": 100, "C2": 5}},
	}

	_, err := service.ComputeSAW(alternatives, criteriaInvalid)
	if err == nil {
		t.Errorf("Expected error for invalid total weight, got nil")
	}
}

// TestComputeSAW_CostAttribute menguji normalisasi kriteria bertipe cost
func TestComputeSAW_CostAttribute(t *testing.T) {
	criteria := []service.SAWCriteria{
		{Code: "C1", AttributeType: "cost", Weight: 0.50},
		{Code: "C2", AttributeType: "benefit", Weight: 0.50},
	}
	alternatives := []service.SAWAlternative{
		{UnitID: 1, UnitCode: "A1", RawValues: map[string]float64{"C1": 10, "C2": 50}},
		{UnitID: 2, UnitCode: "A2", RawValues: map[string]float64{"C1": 20, "C2": 100}},
	}

	result, err := service.ComputeSAW(alternatives, criteria)
	if err != nil {
		t.Fatalf("ComputeSAW cost error: %v", err)
	}

	// Untuk C1 (cost): min C1 = 10
	// A1: r11 = 10/10 = 1.0
	// A2: r21 = 10/20 = 0.5
	for _, alt := range result.Alternatives {
		if alt.UnitCode == "A1" && alt.NormalizedValues["C1"] != 1.0 {
			t.Errorf("Expected A1 normalized C1 = 1.0, got %.4f", alt.NormalizedValues["C1"])
		}
		if alt.UnitCode == "A2" && alt.NormalizedValues["C1"] != 0.5 {
			t.Errorf("Expected A2 normalized C1 = 0.5, got %.4f", alt.NormalizedValues["C1"])
		}
	}
}


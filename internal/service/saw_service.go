package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"erp-anugrah-maha-tunggal/internal/domain"
	"gorm.io/gorm"
)

// ====================================================================================
// ENGINE KOMPUTASI SAW (SIMPLE ADDITIVE WEIGHTING)
// ====================================================================================

// SAWCriteria struct parameter input kriteria untuk algoritma SAW
type SAWCriteria struct {
	ID            uint    `json:"id"`
	Code          string  `json:"code"`           // C1, C2, C3, C4
	Name          string  `json:"name"`           // Nama kriteria
	AttributeType string  `json:"attribute_type"` // benefit / cost
	Weight        float64 `json:"weight"`         // Wj (0.0 s.d 1.0)
	Unit          string  `json:"unit"`           // Jam, Tahun, Kali, Hari
}

// SAWAlternative struct data alternatif unit forklift untuk dinilai
type SAWAlternative struct {
	UnitID      uint               `json:"unit_id"`
	UnitCode    string             `json:"unit_code"`
	Brand       string             `json:"brand"`
	Model       string             `json:"model"`
	CapacityTon float64            `json:"capacity_ton"`
	RawValues   map[string]float64 `json:"raw_values"` // map["C1"]=hour_meter, map["C2"]=machine_age, dst
}

// SAWAlternativeResult struct hasil perhitungan per unit
type SAWAlternativeResult struct {
	UnitID           uint               `json:"unit_id"`
	UnitCode         string             `json:"unit_code"`
	Brand            string             `json:"brand"`
	Model            string             `json:"model"`
	CapacityTon      float64            `json:"capacity_ton"`
	RawValues        map[string]float64 `json:"raw_values"`
	NormalizedValues map[string]float64 `json:"normalized_values"` // Rij
	ScoreV           float64            `json:"score_v"`           // Vi
	RankOrder        int                `json:"rank_order"`        // Peringkat 1..n
	RiskStatus       string             `json:"risk_status"`       // CRITICAL, WARNING, SAFE
	RiskColor        string             `json:"risk_color"`        // merah, kuning, hijau
	Recommendation   string             `json:"recommendation"`
	SPKStatus        string             `json:"spk_status"`        // NONE / ISSUED
}

// SAWCalculationResult struct output lengkap kalkulasi SAW
type SAWCalculationResult struct {
	Criteria      []SAWCriteria          `json:"criteria"`
	Alternatives  []SAWAlternativeResult `json:"alternatives"`
	MaxValues     map[string]float64     `json:"max_values"`
	MinValues     map[string]float64     `json:"min_values"`
	TotalUnits    int                    `json:"total_units"`
	CriticalCount int                    `json:"critical_count"`
	WarningCount  int                    `json:"warning_count"`
	SafeCount     int                    `json:"safe_count"`
	CalculatedAt  time.Time              `json:"calculated_at"`
}

// ComputeSAW menjalankan fungsi komputasi murni algoritma SAW terisolasi
// Parameter:
// - alternatives: Matriks keputusan X (kumpulan baris alternatif dan nilai variabel C1..Cm)
// - criteria: Vektor bobot preferensi W dan sifat atribut (benefit / cost)
func ComputeSAW(alternatives []SAWAlternative, criteria []SAWCriteria) (*SAWCalculationResult, error) {
	if len(alternatives) == 0 {
		return nil, errors.New("alternatif unit forklift tidak boleh kosong")
	}
	if len(criteria) == 0 {
		return nil, errors.New("daftar kriteria SAW tidak boleh kosong")
	}

	// 1. Validasi Total Bobot (Total Wj harus mendekati 1.0 atau 100%)
	var totalWeight float64
	for _, c := range criteria {
		totalWeight += c.Weight
	}
	if math.Abs(totalWeight-1.0) > 0.01 && math.Abs(totalWeight-100.0) > 0.1 {
		return nil, fmt.Errorf("total bobot kriteria harus bernilai 1.0 (100%%), saat ini: %.4f", totalWeight)
	}

	// Normalisasi bobot jika pengguna menginput dalam persentase (e.g. 30 bukan 0.3)
	normalizedCriteria := make([]SAWCriteria, len(criteria))
	copy(normalizedCriteria, criteria)
	if totalWeight > 10.0 {
		for i := range normalizedCriteria {
			normalizedCriteria[i].Weight = normalizedCriteria[i].Weight / 100.0
		}
	}

	// 2. Tentukan Nilai Ekstrem (Max dan Min) per Kriteria
	maxValues := make(map[string]float64)
	minValues := make(map[string]float64)

	for _, c := range normalizedCriteria {
		code := c.Code
		first := true
		for _, alt := range alternatives {
			val := alt.RawValues[code]
			if first {
				maxValues[code] = val
				minValues[code] = val
				first = false
			} else {
				if val > maxValues[code] {
					maxValues[code] = val
				}
				if val < minValues[code] {
					minValues[code] = val
				}
			}
		}
	}

	// 3. Normalisasi Matriks R (rij) & Perkalian Terbobot Nilai Preferensi Vi
	results := make([]SAWAlternativeResult, len(alternatives))
	criticalCount := 0
	warningCount := 0
	safeCount := 0

	for i, alt := range alternatives {
		normMap := make(map[string]float64)
		var scoreV float64

		for _, c := range normalizedCriteria {
			code := c.Code
			xVal := alt.RawValues[code]
			var rVal float64

			if c.AttributeType == "cost" {
				// Atribut Biaya (Cost): rij = min(x_ij) / x_ij
				if xVal > 0 {
					rVal = minValues[code] / xVal
				} else {
					// Penanganan pembagi nol
					rVal = 1.0
				}
			} else {
				// Atribut Keuntungan (Benefit): rij = x_ij / max(x_ij)
				if maxValues[code] > 0 {
					rVal = xVal / maxValues[code]
				} else {
					rVal = 0.0
				}
			}

			// Pembulatan 4 desimal
			rVal = math.Round(rVal*10000) / 10000
			normMap[code] = rVal

			// Vi = SUM( wj * rij )
			scoreV += c.Weight * rVal
		}

		scoreV = math.Round(scoreV*10000) / 10000

		// Klasifikasi Risiko & Rekomendasi
		// Merah (Kritis): Vi >= 0.70 -> Segera Servis
		// Kuning (Waspada): 0.40 <= Vi < 0.70 -> Waspada / Mendekati Batas
		// Hijau (Aman): Vi < 0.40 -> Aman / Kondisi Prima
		var riskStatus, riskColor, recommendation string
		if scoreV >= 0.65 {
			riskStatus = "CRITICAL"
			riskColor = "merah"
			recommendation = "Prioritas Utama Servis: Segera terbitkan SPK bengkel untuk inspeksi menyeluruh sistem mesin dan hidrolik."
			criticalCount++
		} else if scoreV >= 0.40 {
			riskStatus = "WARNING"
			riskColor = "kuning"
			recommendation = "Waspada: Unit mendekati ambang pemeliharaan berkala, jadwalkan pengecekan oli dan filter."
			warningCount++
		} else {
			riskStatus = "SAFE"
			riskColor = "hijau"
			recommendation = "Kondisi Prima: Unit beroperasi optimal dan aman untuk penugasan operasional sewa berikutnya."
			safeCount++
		}

		results[i] = SAWAlternativeResult{
			UnitID:           alt.UnitID,
			UnitCode:         alt.UnitCode,
			Brand:            alt.Brand,
			Model:            alt.Model,
			CapacityTon:      alt.CapacityTon,
			RawValues:        alt.RawValues,
			NormalizedValues: normMap,
			ScoreV:           scoreV,
			RiskStatus:       riskStatus,
			RiskColor:        riskColor,
			Recommendation:   recommendation,
			SPKStatus:        "NONE",
		}
	}

	// 4. Urutkan Array Hasil Peringkat Secara Menurun (Descending Sort by ScoreV)
	sort.Slice(results, func(i, j int) bool {
		return results[i].ScoreV > results[j].ScoreV
	})

	// Beri nomor urut peringkat (RankOrder 1..n)
	for i := range results {
		results[i].RankOrder = i + 1
	}

	return &SAWCalculationResult{
		Criteria:      normalizedCriteria,
		Alternatives:  results,
		MaxValues:     maxValues,
		MinValues:     minValues,
		TotalUnits:    len(results),
		CriticalCount: criticalCount,
		WarningCount:  warningCount,
		SafeCount:     safeCount,
		CalculatedAt:  time.Now(),
	}, nil
}

// ====================================================================================
// SERVICE REPOSITORY & BUSINESS LOGIC
// ====================================================================================

type SAWService interface {
	GetCriteriaAndWeights(ctx context.Context) ([]domain.Kriteria, []domain.BobotKriteria, error)
	UpdateWeights(ctx context.Context, weights []WeightInputRequest, userID *uint) error
	CalculateAndSave(ctx context.Context, userID *uint, notes string) (*SAWCalculationResult, error)
	GetLatestRanking(ctx context.Context) (*SAWCalculationResult, error)
	GetRankingHistory(ctx context.Context, limit int) ([]domain.HasilRanking, error)

	// Log Operasional
	CreateOperationalLog(ctx context.Context, req OperationalLogRequest, userID *uint) (*domain.OperationalLog, error)
	GetOperationalLogs(ctx context.Context, unitID *uint, limit int) ([]domain.OperationalLog, error)

	// SPK Pemeliharaan Bengkel
	CreateSPKMaintenance(ctx context.Context, req SPKMaintenanceRequest, userID *uint) (*domain.SPKMaintenance, error)
	GetSPKMaintenanceList(ctx context.Context) ([]domain.SPKMaintenance, error)
	GetSPKMaintenanceByID(ctx context.Context, id uint) (*domain.SPKMaintenance, error)
}

type sawService struct {
	db *gorm.DB
}

func NewSAWService(db *gorm.DB) SAWService {
	return &sawService{db: db}
}

type WeightInputRequest struct {
	KriteriaID uint    `json:"kriteria_id" binding:"required"`
	Weight     float64 `json:"weight" binding:"required"` // e.g. 0.30 atau 30.0
}

type OperationalLogRequest struct {
	UnitID        uint    `json:"unit_id" binding:"required"`
	LogDate       string  `json:"log_date" binding:"required"`
	HourMeter     float64 `json:"hour_meter" binding:"required"`
	MachineAge    float64 `json:"machine_age" binding:"required"`
	BreakdownFreq float64 `json:"breakdown_freq"`
	IdleDays      float64 `json:"idle_days"`
	Notes         string  `json:"notes"`
}

type SPKMaintenanceRequest struct {
	UnitID         uint    `json:"unit_id" binding:"required"`
	RiskScore      float64 `json:"risk_score"`
	RiskStatus     string  `json:"risk_status"`
	MechanicName   string  `json:"mechanic_name" binding:"required"`
	ComplaintNotes string  `json:"complaint_notes" binding:"required"`
	ScheduledDate  string  `json:"scheduled_date" binding:"required"`
	Notes          string  `json:"notes"`
}

// GetCriteriaAndWeights mengambil master 4 kriteria dan nilai bobot saat ini
func (s *sawService) GetCriteriaAndWeights(ctx context.Context) ([]domain.Kriteria, []domain.BobotKriteria, error) {
	var kriteriaList []domain.Kriteria
	if err := s.db.WithContext(ctx).Order("code ASC").Find(&kriteriaList).Error; err != nil {
		return nil, nil, err
	}

	var bobotList []domain.BobotKriteria
	if err := s.db.WithContext(ctx).Preload("Kriteria").Find(&bobotList).Error; err != nil {
		return nil, nil, err
	}

	return kriteriaList, bobotList, nil
}

// UpdateWeights memperbarui bobot kriteria dan memvalidasi total 100%
func (s *sawService) UpdateWeights(ctx context.Context, weights []WeightInputRequest, userID *uint) error {
	if len(weights) == 0 {
		return errors.New("bobot kriteria tidak boleh kosong")
	}

	var sum float64
	for _, w := range weights {
		sum += w.Weight
	}

	// Normalisasi desimal jika dikirim dalam persen (misal 30 -> 0.3)
	isPercent := sum > 10.0
	if isPercent {
		if math.Abs(sum-100.0) > 0.1 {
			return fmt.Errorf("total persentase bobot harus 100%% (saat ini: %.2f%%)", sum)
		}
	} else {
		if math.Abs(sum-1.0) > 0.01 {
			return fmt.Errorf("total bobot harus bernilai 1.0 (saat ini: %.4f)", sum)
		}
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, w := range weights {
			weightVal := w.Weight
			percentageVal := w.Weight
			if isPercent {
				weightVal = w.Weight / 100.0
			} else {
				percentageVal = w.Weight * 100.0
			}

			var bobot domain.BobotKriteria
			if err := tx.Where("kriteria_id = ?", w.KriteriaID).First(&bobot).Error; err != nil {
				// Buat baru jika belum ada
				bobot = domain.BobotKriteria{
					KriteriaID: w.KriteriaID,
					Weight:     weightVal,
					Percentage: percentageVal,
					SetBy:      userID,
					Notes:      fmt.Sprintf("Diperbarui oleh manajemen sebesar %.2f%%", percentageVal),
				}
				if err := tx.Create(&bobot).Error; err != nil {
					return err
				}
			} else {
				// Update
				bobot.Weight = weightVal
				bobot.Percentage = percentageVal
				bobot.SetBy = userID
				bobot.Notes = fmt.Sprintf("Diperbarui oleh manajemen sebesar %.2f%%", percentageVal)
				if err := tx.Save(&bobot).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// CalculateAndSave menjalankan komputasi SAW dan menyimpan snapshot audit trail
func (s *sawService) CalculateAndSave(ctx context.Context, userID *uint, notes string) (*SAWCalculationResult, error) {
	// 1. Ambil kriteria dan bobot
	var bobotList []domain.BobotKriteria
	if err := s.db.WithContext(ctx).Preload("Kriteria").Find(&bobotList).Error; err != nil {
		return nil, err
	}
	if len(bobotList) == 0 {
		return nil, errors.New("bobot kriteria belum diatur di database")
	}

	criteria := make([]SAWCriteria, len(bobotList))
	for i, b := range bobotList {
		criteria[i] = SAWCriteria{
			ID:            b.KriteriaID,
			Code:          b.Kriteria.Code,
			Name:          b.Kriteria.Name,
			AttributeType: b.Kriteria.AttributeType,
			Weight:        b.Weight,
			Unit:          b.Kriteria.Unit,
		}
	}

	// 2. Ambil data unit armada
	var units []domain.ForkliftUnit
	if err := s.db.WithContext(ctx).Find(&units).Error; err != nil {
		return nil, err
	}
	if len(units) == 0 {
		return nil, errors.New("tidak ada unit forklift terdaftar untuk dinilai")
	}

	// Bentuk alternatif matriks X
	alternatives := make([]SAWAlternative, len(units))
	for i, u := range units {
		raw := map[string]float64{
			"C1": u.HourMeter,
			"C2": u.MachineAge,
			"C3": u.BreakdownFrequency,
			"C4": u.IdleDays,
		}
		alternatives[i] = SAWAlternative{
			UnitID:      u.ID,
			UnitCode:    u.UnitCode,
			Brand:       u.Brand,
			Model:       u.Model,
			CapacityTon: u.CapacityTon,
			RawValues:   raw,
		}
	}

	// 3. Eksekusi Engine SAW
	calcResult, err := ComputeSAW(alternatives, criteria)
	if err != nil {
		return nil, err
	}

	// 4. Simpan Snapshot ke tabel hasil_ranking & detail_ranking
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		code := fmt.Sprintf("SAW-%s-%03d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000)
		if notes == "" {
			notes = fmt.Sprintf("Kalkulasi rekomendasi prioritas pemeliharaan per %s", time.Now().Format("02-01-2006 15:04 WIB"))
		}

		snapshot := domain.HasilRanking{
			CalculationCode: code,
			CalculatedAt:    calcResult.CalculatedAt,
			CalculatedBy:    userID,
			TotalUnits:      calcResult.TotalUnits,
			CriticalCount:   calcResult.CriticalCount,
			WarningCount:    calcResult.WarningCount,
			SafeCount:       calcResult.SafeCount,
			Notes:           notes,
		}
		if err := tx.Create(&snapshot).Error; err != nil {
			return err
		}

		// Simpan detail per unit dan perbarui status unit di forklift_units
		for _, alt := range calcResult.Alternatives {
			detail := domain.DetailRanking{
				HasilRankingID: snapshot.ID,
				UnitID:         alt.UnitID,
				RankOrder:      alt.RankOrder,
				C1Raw:          alt.RawValues["C1"],
				C2Raw:          alt.RawValues["C2"],
				C3Raw:          alt.RawValues["C3"],
				C4Raw:          alt.RawValues["C4"],
				C1Normalized:   alt.NormalizedValues["C1"],
				C2Normalized:   alt.NormalizedValues["C2"],
				C3Normalized:   alt.NormalizedValues["C3"],
				C4Normalized:   alt.NormalizedValues["C4"],
				ScoreV:         alt.ScoreV,
				RiskStatus:     alt.RiskStatus,
				RiskColor:      alt.RiskColor,
				Recommendation: alt.Recommendation,
				SPKStatus:      "NONE",
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			// Update skor dan status risiko pada master forklift_units
			if err := tx.Model(&domain.ForkliftUnit{}).Where("id = ?", alt.UnitID).Updates(map[string]interface{}{
				"last_calculated_score": alt.ScoreV,
				"last_risk_status":      alt.RiskStatus,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan snapshot ranking: %w", err)
	}

	return calcResult, nil
}

// GetLatestRanking mengambil ranking rekomendasi terkini
func (s *sawService) GetLatestRanking(ctx context.Context) (*SAWCalculationResult, error) {
	var latestSnapshot domain.HasilRanking
	err := s.db.WithContext(ctx).Order("calculated_at DESC").First(&latestSnapshot).Error
	if err != nil {
		// Jika belum ada snapshot, hitung langsung
		return s.CalculateAndSave(ctx, nil, "Kalkulasi inisialisasi awal")
	}

	var details []domain.DetailRanking
	if err := s.db.WithContext(ctx).
		Preload("Unit").
		Where("hasil_ranking_id = ?", latestSnapshot.ID).
		Order("rank_order ASC").
		Find(&details).Error; err != nil {
		return nil, err
	}

	// Periksa apakah unit sudah punya SPK aktif
	var activeSPKUnitIDs []uint
	s.db.WithContext(ctx).Model(&domain.SPKMaintenance{}).
		Where("status IN ('DRAFT', 'IN_PROGRESS')").
		Pluck("unit_id", &activeSPKUnitIDs)
	activeSPKMap := make(map[uint]bool)
	for _, id := range activeSPKUnitIDs {
		activeSPKMap[id] = true
	}

	var bobotList []domain.BobotKriteria
	s.db.WithContext(ctx).Preload("Kriteria").Find(&bobotList)
	criteria := make([]SAWCriteria, len(bobotList))
	for i, b := range bobotList {
		criteria[i] = SAWCriteria{
			ID:            b.KriteriaID,
			Code:          b.Kriteria.Code,
			Name:          b.Kriteria.Name,
			AttributeType: b.Kriteria.AttributeType,
			Weight:        b.Weight,
			Unit:          b.Kriteria.Unit,
		}
	}

	alts := make([]SAWAlternativeResult, len(details))
	for i, d := range details {
		spkStatus := "NONE"
		if activeSPKMap[d.UnitID] {
			spkStatus = "ISSUED"
		}

		unitCode := ""
		brand := ""
		model := ""
		capTon := 0.0
		if d.Unit != nil {
			unitCode = d.Unit.UnitCode
			brand = d.Unit.Brand
			model = d.Unit.Model
			capTon = d.Unit.CapacityTon
		}

		alts[i] = SAWAlternativeResult{
			UnitID:      d.UnitID,
			UnitCode:    unitCode,
			Brand:       brand,
			Model:       model,
			CapacityTon: capTon,
			RawValues: map[string]float64{
				"C1": d.C1Raw,
				"C2": d.C2Raw,
				"C3": d.C3Raw,
				"C4": d.C4Raw,
			},
			NormalizedValues: map[string]float64{
				"C1": d.C1Normalized,
				"C2": d.C2Normalized,
				"C3": d.C3Normalized,
				"C4": d.C4Normalized,
			},
			ScoreV:         d.ScoreV,
			RankOrder:      d.RankOrder,
			RiskStatus:     d.RiskStatus,
			RiskColor:      d.RiskColor,
			Recommendation: d.Recommendation,
			SPKStatus:      spkStatus,
		}
	}

	return &SAWCalculationResult{
		Criteria:      criteria,
		Alternatives:  alts,
		TotalUnits:    latestSnapshot.TotalUnits,
		CriticalCount: latestSnapshot.CriticalCount,
		WarningCount:  latestSnapshot.WarningCount,
		SafeCount:     latestSnapshot.SafeCount,
		CalculatedAt:  latestSnapshot.CalculatedAt,
	}, nil
}

// GetRankingHistory mengambil riwayat snapshot hasil perhitungan SAW
func (s *sawService) GetRankingHistory(ctx context.Context, limit int) ([]domain.HasilRanking, error) {
	if limit <= 0 {
		limit = 10
	}
	var history []domain.HasilRanking
	err := s.db.WithContext(ctx).
		Preload("Calculator").
		Order("calculated_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

// CreateOperationalLog mencatat variabel mentah harian dan otomatis mentrigger update unit & SAW
func (s *sawService) CreateOperationalLog(ctx context.Context, req OperationalLogRequest, userID *uint) (*domain.OperationalLog, error) {
	var unit domain.ForkliftUnit
	if err := s.db.WithContext(ctx).First(&unit, req.UnitID).Error; err != nil {
		return nil, fmt.Errorf("unit forklift dengan ID %d tidak ditemukan", req.UnitID)
	}

	logEntry := domain.OperationalLog{
		UnitID:        req.UnitID,
		LogDate:       req.LogDate,
		HourMeter:     req.HourMeter,
		MachineAge:    req.MachineAge,
		BreakdownFreq: req.BreakdownFreq,
		IdleDays:      req.IdleDays,
		RecordedBy:    userID,
		Notes:         req.Notes,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&logEntry).Error; err != nil {
			return err
		}

		// Update variabel mentah C1..C4 di master forklift_units
		return tx.Model(&domain.ForkliftUnit{}).Where("id = ?", req.UnitID).Updates(map[string]interface{}{
			"hour_meter":          req.HourMeter,
			"machine_age":         req.MachineAge,
			"breakdown_frequency": req.BreakdownFreq,
			"idle_days":           req.IdleDays,
		}).Error
	})

	if err != nil {
		return nil, err
	}

	// Trigger hitung ulang ranking otomatis di background/sync agar data ranking selalu up to date
	go func() {
		_, _ = s.CalculateAndSave(context.Background(), userID, fmt.Sprintf("Kalkulasi otomatis pasca-input log operasional unit %s", unit.UnitCode))
	}()

	s.db.Preload("Unit").Preload("Recorder").First(&logEntry, logEntry.ID)
	return &logEntry, nil
}

// GetOperationalLogs mengambil daftar log operasional
func (s *sawService) GetOperationalLogs(ctx context.Context, unitID *uint, limit int) ([]domain.OperationalLog, error) {
	if limit <= 0 {
		limit = 50
	}
	query := s.db.WithContext(ctx).Preload("Unit").Preload("Recorder").Order("log_date DESC, created_at DESC")
	if unitID != nil && *unitID > 0 {
		query = query.Where("unit_id = ?", *unitID)
	}
	var logs []domain.OperationalLog
	err := query.Limit(limit).Find(&logs).Error
	return logs, err
}

// CreateSPKMaintenance menerbitkan SPK Pemeliharaan, mengubah status unit menjadi 'MAINTENANCE' di ERP
func (s *sawService) CreateSPKMaintenance(ctx context.Context, req SPKMaintenanceRequest, userID *uint) (*domain.SPKMaintenance, error) {
	var unit domain.ForkliftUnit
	if err := s.db.WithContext(ctx).First(&unit, req.UnitID).Error; err != nil {
		return nil, fmt.Errorf("unit forklift dengan ID %d tidak ditemukan", req.UnitID)
	}

	spkNumber := fmt.Sprintf("SPK-MT-%s-%03d", time.Now().Format("20060102"), (time.Now().Unix()%900)+100)

	spk := domain.SPKMaintenance{
		SPKNumber:      spkNumber,
		UnitID:         req.UnitID,
		RiskScore:      req.RiskScore,
		RiskStatus:     req.RiskStatus,
		MechanicName:   req.MechanicName,
		ComplaintNotes: req.ComplaintNotes,
		ScheduledDate:  req.ScheduledDate,
		Status:         "IN_PROGRESS",
		IssuedBy:       userID,
		IssuedAt:       time.Now(),
		Notes:          req.Notes,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&spk).Error; err != nil {
			return err
		}

		// Update status unit di ERP menjadi MAINTENANCE ("Dalam Perawatan")
		oldStatus := unit.Status
		if err := tx.Model(&domain.ForkliftUnit{}).Where("id = ?", req.UnitID).Update("status", "MAINTENANCE").Error; err != nil {
			return err
		}

		// Catat ke UnitMovementLog (Audit Trail status)
		movementLog := domain.UnitMovementLog{
			UnitID:     req.UnitID,
			StatusFrom: oldStatus,
			StatusTo:   "MAINTENANCE",
			Notes:      fmt.Sprintf("Penerbitan Surat Perintah Kerja (SPK) Bengkel No: %s. Keluhan: %s", spkNumber, req.ComplaintNotes),
			LoggedAt:   time.Now(),
		}
		return tx.Create(&movementLog).Error
	})

	if err != nil {
		return nil, err
	}

	s.db.Preload("Unit").Preload("Issuer").First(&spk, spk.ID)
	return &spk, nil
}

// GetSPKMaintenanceList mengambil daftar SPK pemeliharaan
func (s *sawService) GetSPKMaintenanceList(ctx context.Context) ([]domain.SPKMaintenance, error) {
	var list []domain.SPKMaintenance
	err := s.db.WithContext(ctx).
		Preload("Unit").
		Preload("Issuer").
		Order("issued_at DESC").
		Find(&list).Error
	return list, err
}

// GetSPKMaintenanceByID mengambil satu SPK untuk dicetak fisik
func (s *sawService) GetSPKMaintenanceByID(ctx context.Context, id uint) (*domain.SPKMaintenance, error) {
	var spk domain.SPKMaintenance
	err := s.db.WithContext(ctx).
		Preload("Unit").
		Preload("Issuer").
		First(&spk, id).Error
	return &spk, err
}


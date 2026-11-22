package domain

import (
	"time"
)

type UserRole string

const (
	RoleDirektur    UserRole = "direktur"
	RoleAdminStaff  UserRole = "admin_staff"
	RoleOperasional UserRole = "operasional"
	RoleFinance     UserRole = "finance"
)

// User merepresentasikan data akun sistem ERP
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Username     string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole  `gorm:"type:enum('direktur','admin_staff','operasional','finance');default:'admin_staff'" json:"role"`
	Phone        string    `gorm:"type:varchar(20)" json:"phone"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Customer merepresentasikan klien penyewa (PT, CV, Perorangan)
type Customer struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerCode string    `gorm:"type:varchar(20);unique;not null" json:"customer_code"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Type         string    `gorm:"type:enum('PT','CV','Perorangan','Instansi');default:'PT'" json:"type"`
	PicName      string    `gorm:"type:varchar(100);not null" json:"pic_name"`
	Phone        string    `gorm:"type:varchar(25);not null" json:"phone"`
	Email        string    `gorm:"type:varchar(100)" json:"email"`
	Address      string    `gorm:"type:text;not null" json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ForkliftUnit merepresentasikan aset alat berat / unit forklift
type ForkliftUnit struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UnitCode            string    `gorm:"type:varchar(30);unique;not null" json:"unit_code"`
	Brand               string    `gorm:"type:varchar(50);not null" json:"brand"`
	Model               string    `gorm:"type:varchar(50)" json:"model"`
	CapacityTon         float64   `gorm:"type:decimal(4,1);not null" json:"capacity_ton"`
	FuelType            string    `gorm:"type:enum('Diesel','Electric','Gasoline','LPG');default:'Diesel'" json:"fuel_type"`
	ManufactureYear     int       `json:"manufacture_year"`
	HourlyRate          float64   `gorm:"type:decimal(12,2);default:0" json:"hourly_rate"`
	DailyRate           float64   `gorm:"type:decimal(12,2);default:0" json:"daily_rate"`
	MonthlyRate         float64   `gorm:"type:decimal(12,2);default:0" json:"monthly_rate"`
	Status              string    `gorm:"type:enum('AVAILABLE','RENTED','MAINTENANCE');default:'AVAILABLE'" json:"status"`
	Notes               string    `gorm:"type:text" json:"notes"`
	// Kolom SPK Maintenance SAW
	HourMeter           float64   `gorm:"type:decimal(10,2);default:0" json:"hour_meter"`           // C1 (Jam)
	MachineAge          float64   `gorm:"type:decimal(5,2);default:0" json:"machine_age"`          // C2 (Tahun)
	BreakdownFrequency  float64   `gorm:"type:decimal(5,2);default:0" json:"breakdown_frequency"`  // C3 (Kali)
	IdleDays            float64   `gorm:"type:decimal(5,2);default:0" json:"idle_days"`            // C4 (Hari)
	LastCalculatedScore float64   `gorm:"type:decimal(8,4);default:0" json:"last_calculated_score"` // Skor Vi
	LastRiskStatus      string    `gorm:"type:varchar(20);default:'SAFE'" json:"last_risk_status"`  // CRITICAL / WARNING / SAFE
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Operator merepresentasikan driver / teknisi operator forklift
type Operator struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nip       string    `gorm:"type:varchar(30);unique" json:"nip"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Phone     string    `gorm:"type:varchar(25);not null" json:"phone"`
	SioNumber string    `gorm:"type:varchar(50)" json:"sio_number"` // Surat Izin Operator
	Status    string    `gorm:"type:enum('STANDBY','ON_DUTY','OFF');default:'STANDBY'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RentalOrder merepresentasikan order sewa dari customer (PO Customer)
type RentalOrder struct {
	ID                  uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber         string        `gorm:"type:varchar(40);unique;not null" json:"order_number"`
	CustomerID          uint          `gorm:"not null" json:"customer_id"`
	Customer            *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	OrderDate           string        `gorm:"type:date;not null" json:"order_date"`
	StartDate           string        `gorm:"type:date;not null" json:"start_date"`
	EndDate             *string       `gorm:"type:date" json:"end_date"`
	RentalDurationType  string        `gorm:"type:enum('JAM','HARI','BULAN');default:'HARI'" json:"rental_duration_type"`
	DurationValue       int           `gorm:"default:1" json:"duration_value"`
	ProjectLocation     string        `gorm:"type:text;not null" json:"project_location"`
	LocationPicName     string        `gorm:"type:varchar(100)" json:"location_pic_name"`
	LocationPicPhone    string        `gorm:"type:varchar(25)" json:"location_pic_phone"`
	RequiredCapacityTon float64       `gorm:"type:decimal(4,1);not null" json:"required_capacity_ton"`
	UnitID              *uint         `json:"unit_id"`
	Unit                *ForkliftUnit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	AgreedPrice         float64       `gorm:"type:decimal(12,2);default:0" json:"agreed_price"`
	Status              string        `gorm:"type:enum('PENDING','APPROVED','ACTIVE','COMPLETED','CANCELLED');default:'PENDING'" json:"status"`
	PoReference         string        `gorm:"type:varchar(50)" json:"po_reference"`
	Notes               string        `gorm:"type:text" json:"notes"`
	CreatedBy           *uint         `json:"created_by"`
	Creator             *User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

// DeliveryLetter merepresentasikan Dokumen Surat Jalan Resmi untuk operasional
type DeliveryLetter struct {
	ID                uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	LetterNumber      string        `gorm:"type:varchar(40);unique;not null" json:"letter_number"`
	RentalOrderID     uint          `gorm:"not null" json:"rental_order_id"`
	RentalOrder       *RentalOrder  `gorm:"foreignKey:RentalOrderID" json:"rental_order,omitempty"`
	UnitID            uint          `gorm:"not null" json:"unit_id"`
	Unit              *ForkliftUnit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	OperatorID        uint          `gorm:"not null" json:"operator_id"`
	Operator          *Operator     `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	IssueDate         string        `gorm:"type:date;not null" json:"issue_date"`
	DepartureTime     string        `gorm:"type:varchar(50)" json:"departure_time"`
	JobDescription    string        `gorm:"type:text" json:"job_description"`
	RecipientName     string        `gorm:"type:varchar(100)" json:"recipient_name"`
	OperationalStatus string        `gorm:"type:enum('ASSIGNED','ON_THE_WAY','WORKING','FINISHED');default:'ASSIGNED'" json:"operational_status"`
	Notes             string        `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// UnitMovementLog untuk riwayat audit trail pergerakan unit forklift
type UnitMovementLog struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UnitID           uint      `gorm:"not null" json:"unit_id"`
	DeliveryLetterID *uint     `json:"delivery_letter_id"`
	StatusFrom       string    `gorm:"type:varchar(50)" json:"status_from"`
	StatusTo         string    `gorm:"type:varchar(50);not null" json:"status_to"`
	Notes            string    `gorm:"type:text" json:"notes"`
	LoggedAt         time.Time `json:"logged_at"`
}

// Invoice & Pembayaran
type Invoice struct {
	ID            uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceNumber string       `gorm:"type:varchar(40);unique;not null" json:"invoice_number"`
	RentalOrderID uint         `gorm:"not null" json:"rental_order_id"`
	RentalOrder   *RentalOrder `gorm:"foreignKey:RentalOrderID" json:"rental_order,omitempty"`
	TotalAmount   float64      `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	PaidAmount    float64      `gorm:"type:decimal(12,2);default:0" json:"paid_amount"`
	PaymentStatus string       `gorm:"type:enum('UNPAID','PARTIAL','PAID');default:'UNPAID'" json:"payment_status"`
	DueDate       *string      `gorm:"type:date" json:"due_date"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// ====================================================================================
// MODUL SPK PEMELIHARAAN FORKLIFT METODE SAW (SIMPLE ADDITIVE WEIGHTING)
// ====================================================================================

// Kriteria merepresentasikan master kriteria SPK (C1..C4)
type Kriteria struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code          string    `gorm:"type:varchar(10);unique;not null" json:"code"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	AttributeType string    `gorm:"type:enum('benefit','cost');default:'benefit'" json:"attribute_type"` // benefit atau cost
	Unit          string    `gorm:"type:varchar(20)" json:"unit"`                                        // Jam, Tahun, Kali, Hari
	Description   string    `gorm:"type:text" json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Kriteria) TableName() string {
	return "kriteria"
}

// BobotKriteria merepresentasikan bobot preferensi manajemen (W1..W4)
type BobotKriteria struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	KriteriaID uint      `gorm:"not null" json:"kriteria_id"`
	Kriteria   *Kriteria `gorm:"foreignKey:KriteriaID" json:"kriteria,omitempty"`
	Weight     float64   `gorm:"type:decimal(6,4);not null" json:"weight"`      // e.g. 0.3000
	Percentage float64   `gorm:"type:decimal(5,2);default:0" json:"percentage"` // e.g. 30.00
	SetBy      *uint     `json:"set_by"`
	User       *User     `gorm:"foreignKey:SetBy" json:"user,omitempty"`
	Notes      string    `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (BobotKriteria) TableName() string {
	return "bobot_kriteria"
}

// OperationalLog mencatat log variabel mentah armada forklift
type OperationalLog struct {
	ID            uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UnitID        uint          `gorm:"not null" json:"unit_id"`
	Unit          *ForkliftUnit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	LogDate       string        `gorm:"type:date;not null" json:"log_date"`
	HourMeter     float64       `gorm:"type:decimal(10,2);not null" json:"hour_meter"`     // C1 (Jam)
	MachineAge    float64       `gorm:"type:decimal(5,2);not null" json:"machine_age"`     // C2 (Tahun)
	BreakdownFreq float64       `gorm:"type:decimal(5,2);not null" json:"breakdown_freq"`  // C3 (Kali)
	IdleDays      float64       `gorm:"type:decimal(5,2);not null" json:"idle_days"`       // C4 (Hari)
	RecordedBy    *uint         `json:"recorded_by"`
	Recorder      *User         `gorm:"foreignKey:RecordedBy" json:"recorder,omitempty"`
	Notes         string        `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (OperationalLog) TableName() string {
	return "log_operasional"
}

// HasilRanking snapshot audit trail riwayat kalkulasi SAW
type HasilRanking struct {
	ID              uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	CalculationCode string          `gorm:"type:varchar(50);unique;not null" json:"calculation_code"`
	CalculatedAt    time.Time       `json:"calculated_at"`
	CalculatedBy    *uint           `json:"calculated_by"`
	Calculator      *User           `gorm:"foreignKey:CalculatedBy" json:"calculator,omitempty"`
	TotalUnits      int             `json:"total_units"`
	CriticalCount   int             `json:"critical_count"` // Merah (Kritis)
	WarningCount    int             `json:"warning_count"`  // Kuning (Waspada)
	SafeCount       int             `json:"safe_count"`     // Hijau (Aman)
	Notes           string          `gorm:"type:text" json:"notes"`
	Details         []DetailRanking `gorm:"foreignKey:HasilRankingID" json:"details,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

func (HasilRanking) TableName() string {
	return "hasil_ranking"
}

// DetailRanking rincian peringkat per unit forklift dalam satu snapshot kalkulasi
type DetailRanking struct {
	ID             uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	HasilRankingID uint          `gorm:"not null" json:"hasil_ranking_id"`
	UnitID         uint          `gorm:"not null" json:"unit_id"`
	Unit           *ForkliftUnit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	RankOrder      int           `gorm:"not null" json:"rank_order"`
	C1Raw          float64       `gorm:"type:decimal(10,2)" json:"c1_raw"`
	C2Raw          float64       `gorm:"type:decimal(5,2)" json:"c2_raw"`
	C3Raw          float64       `gorm:"type:decimal(5,2)" json:"c3_raw"`
	C4Raw          float64       `gorm:"type:decimal(5,2)" json:"c4_raw"`
	C1Normalized   float64       `gorm:"type:decimal(8,4)" json:"c1_normalized"`
	C2Normalized   float64       `gorm:"type:decimal(8,4)" json:"c2_normalized"`
	C3Normalized   float64       `gorm:"type:decimal(8,4)" json:"c3_normalized"`
	C4Normalized   float64       `gorm:"type:decimal(8,4)" json:"c4_normalized"`
	ScoreV         float64       `gorm:"type:decimal(8,4)" json:"score_v"` // Nilai Vi
	RiskStatus     string        `gorm:"type:varchar(20)" json:"risk_status"` // CRITICAL, WARNING, SAFE
	RiskColor      string        `gorm:"type:varchar(20)" json:"risk_color"`  // merah, kuning, hijau
	Recommendation string        `gorm:"type:text" json:"recommendation"`
	SPKStatus      string        `gorm:"type:varchar(20);default:'NONE'" json:"spk_status"` // NONE, ISSUED
}

func (DetailRanking) TableName() string {
	return "detail_ranking"
}

// SPKMaintenance Surat Perintah Kerja (SPK) Pemeliharaan Bengkel
type SPKMaintenance struct {
	ID             uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	SPKNumber      string        `gorm:"type:varchar(50);unique;not null" json:"spk_number"`
	UnitID         uint          `gorm:"not null" json:"unit_id"`
	Unit           *ForkliftUnit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	RiskScore      float64       `gorm:"type:decimal(8,4)" json:"risk_score"`
	RiskStatus     string        `gorm:"type:varchar(20)" json:"risk_status"`
	MechanicName   string        `gorm:"type:varchar(100);not null" json:"mechanic_name"`
	ComplaintNotes string        `gorm:"type:text;not null" json:"complaint_notes"`
	ScheduledDate  string        `gorm:"type:date;not null" json:"scheduled_date"`
	Status         string        `gorm:"type:enum('DRAFT','IN_PROGRESS','COMPLETED','CANCELLED');default:'IN_PROGRESS'" json:"status"`
	IssuedBy       *uint         `json:"issued_by"`
	Issuer         *User         `gorm:"foreignKey:IssuedBy" json:"issuer,omitempty"`
	IssuedAt       time.Time     `json:"issued_at"`
	Notes          string        `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func (SPKMaintenance) TableName() string {
	return "spk_maintenance"
}

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
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UnitCode        string    `gorm:"type:varchar(30);unique;not null" json:"unit_code"`
	Brand           string    `gorm:"type:varchar(50);not null" json:"brand"`
	Model           string    `gorm:"type:varchar(50)" json:"model"`
	CapacityTon     float64   `gorm:"type:decimal(4,1);not null" json:"capacity_ton"`
	FuelType        string    `gorm:"type:enum('Diesel','Electric','Gasoline','LPG');default:'Diesel'" json:"fuel_type"`
	ManufactureYear int       `json:"manufacture_year"`
	HourlyRate      float64   `gorm:"type:decimal(12,2);default:0" json:"hourly_rate"`
	DailyRate       float64   `gorm:"type:decimal(12,2);default:0" json:"daily_rate"`
	MonthlyRate     float64   `gorm:"type:decimal(12,2);default:0" json:"monthly_rate"`
	Status          string    `gorm:"type:enum('AVAILABLE','RENTED','MAINTENANCE');default:'AVAILABLE'" json:"status"`
	Notes           string    `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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

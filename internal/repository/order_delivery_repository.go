package repository

import (
	"erp-anugrah-maha-tunggal/internal/domain"

	"gorm.io/gorm"
)

type RentalOrderRepository interface {
	FindAll(status string) ([]domain.RentalOrder, error)
	FindByID(id uint) (*domain.RentalOrder, error)
	Create(order *domain.RentalOrder) error
	Update(order *domain.RentalOrder) error
	UpdateStatus(id uint, status string) error
	Count() (int64, error)
	CountActive() (int64, error)
}

type rentalOrderRepository struct {
	db *gorm.DB
}

func NewRentalOrderRepository(db *gorm.DB) RentalOrderRepository {
	return &rentalOrderRepository{db: db}
}

func (r *rentalOrderRepository) FindAll(status string) ([]domain.RentalOrder, error) {
	var orders []domain.RentalOrder
	q := r.db.Preload("Customer").Preload("Unit").Preload("Creator").Order("id desc")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&orders).Error
	return orders, err
}

func (r *rentalOrderRepository) FindByID(id uint) (*domain.RentalOrder, error) {
	var order domain.RentalOrder
	err := r.db.Preload("Customer").Preload("Unit").Preload("Creator").First(&order, id).Error
	return &order, err
}

func (r *rentalOrderRepository) Create(order *domain.RentalOrder) error {
	return r.db.Create(order).Error
}

func (r *rentalOrderRepository) Update(order *domain.RentalOrder) error {
	return r.db.Save(order).Error
}

func (r *rentalOrderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.RentalOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *rentalOrderRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.RentalOrder{}).Count(&count).Error
	return count, err
}

func (r *rentalOrderRepository) CountActive() (int64, error) {
	var count int64
	err := r.db.Model(&domain.RentalOrder{}).Where("status IN ('APPROVED', 'ACTIVE')").Count(&count).Error
	return count, err
}

// ================= DELIVERY LETTER REPOSITORY =================
type DeliveryLetterRepository interface {
	FindAll(status string) ([]domain.DeliveryLetter, error)
	FindByID(id uint) (*domain.DeliveryLetter, error)
	FindByOrderID(orderID uint) ([]domain.DeliveryLetter, error)
	Create(letter *domain.DeliveryLetter) error
	Update(letter *domain.DeliveryLetter) error
	UpdateOperationalStatus(id uint, status string) error
	CreateMovementLog(log *domain.UnitMovementLog) error
	GetMovementLogs(unitID uint) ([]domain.UnitMovementLog, error)
}

type deliveryLetterRepository struct {
	db *gorm.DB
}

func NewDeliveryLetterRepository(db *gorm.DB) DeliveryLetterRepository {
	return &deliveryLetterRepository{db: db}
}

func (r *deliveryLetterRepository) FindAll(status string) ([]domain.DeliveryLetter, error) {
	var letters []domain.DeliveryLetter
	q := r.db.Preload("RentalOrder").Preload("RentalOrder.Customer").Preload("Unit").Preload("Operator").Order("id desc")
	if status != "" {
		q = q.Where("operational_status = ?", status)
	}
	err := q.Find(&letters).Error
	return letters, err
}

func (r *deliveryLetterRepository) FindByID(id uint) (*domain.DeliveryLetter, error) {
	var letter domain.DeliveryLetter
	err := r.db.Preload("RentalOrder").Preload("RentalOrder.Customer").Preload("Unit").Preload("Operator").First(&letter, id).Error
	return &letter, err
}

func (r *deliveryLetterRepository) FindByOrderID(orderID uint) ([]domain.DeliveryLetter, error) {
	var letters []domain.DeliveryLetter
	err := r.db.Preload("RentalOrder").Preload("Unit").Preload("Operator").Where("rental_order_id = ?", orderID).Find(&letters).Error
	return letters, err
}

func (r *deliveryLetterRepository) Create(letter *domain.DeliveryLetter) error {
	return r.db.Create(letter).Error
}

func (r *deliveryLetterRepository) Update(letter *domain.DeliveryLetter) error {
	return r.db.Save(letter).Error
}

func (r *deliveryLetterRepository) UpdateOperationalStatus(id uint, status string) error {
	return r.db.Model(&domain.DeliveryLetter{}).Where("id = ?", id).Update("operational_status", status).Error
}

func (r *deliveryLetterRepository) CreateMovementLog(log *domain.UnitMovementLog) error {
	return r.db.Create(log).Error
}

func (r *deliveryLetterRepository) GetMovementLogs(unitID uint) ([]domain.UnitMovementLog, error) {
	var logs []domain.UnitMovementLog
	err := r.db.Where("unit_id = ?", unitID).Order("logged_at desc").Find(&logs).Error
	return logs, err
}

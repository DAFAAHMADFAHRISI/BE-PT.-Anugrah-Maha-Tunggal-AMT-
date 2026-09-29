package repository

import (
	"erp-anugrah-maha-tunggal/internal/domain"

	"gorm.io/gorm"
)

// ================= CUSTOMER REPOSITORY =================
type CustomerRepository interface {
	FindAll() ([]domain.Customer, error)
	FindByID(id uint) (*domain.Customer, error)
	Create(customer *domain.Customer) error
	Update(customer *domain.Customer) error
	Delete(id uint) error
	Count() (int64, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) FindAll() ([]domain.Customer, error) {
	var customers []domain.Customer
	err := r.db.Order("id desc").Find(&customers).Error
	return customers, err
}

func (r *customerRepository) FindByID(id uint) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.First(&customer, id).Error
	return &customer, err
}

func (r *customerRepository) Create(customer *domain.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *domain.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Customer{}, id).Error
}

func (r *customerRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Customer{}).Count(&count).Error
	return count, err
}

// ================= FORKLIFT UNIT REPOSITORY =================
type UnitRepository interface {
	FindAll(statusFilter string) ([]domain.ForkliftUnit, error)
	FindByID(id uint) (*domain.ForkliftUnit, error)
	Create(unit *domain.ForkliftUnit) error
	Update(unit *domain.ForkliftUnit) error
	UpdateStatus(id uint, status string) error
	Delete(id uint) error
	GetStatistics() (total int64, available int64, rented int64, maintenance int64, err error)
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) FindAll(statusFilter string) ([]domain.ForkliftUnit, error) {
	var units []domain.ForkliftUnit
	q := r.db.Order("capacity_ton asc")
	if statusFilter != "" {
		q = q.Where("status = ?", statusFilter)
	}
	err := q.Find(&units).Error
	return units, err
}

func (r *unitRepository) FindByID(id uint) (*domain.ForkliftUnit, error) {
	var unit domain.ForkliftUnit
	err := r.db.First(&unit, id).Error
	return &unit, err
}

func (r *unitRepository) Create(unit *domain.ForkliftUnit) error {
	return r.db.Create(unit).Error
}

func (r *unitRepository) Update(unit *domain.ForkliftUnit) error {
	return r.db.Save(unit).Error
}

func (r *unitRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.ForkliftUnit{}).Where("id = ?", id).Update("status", status).Error
}

func (r *unitRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ForkliftUnit{}, id).Error
}

func (r *unitRepository) GetStatistics() (total int64, available int64, rented int64, maintenance int64, err error) {
	err = r.db.Model(&domain.ForkliftUnit{}).Count(&total).Error
	if err != nil {
		return
	}
	r.db.Model(&domain.ForkliftUnit{}).Where("status = ?", "AVAILABLE").Count(&available)
	r.db.Model(&domain.ForkliftUnit{}).Where("status = ?", "RENTED").Count(&rented)
	r.db.Model(&domain.ForkliftUnit{}).Where("status = ?", "MAINTENANCE").Count(&maintenance)
	return
}

// ================= OPERATOR REPOSITORY =================
type OperatorRepository interface {
	FindAll() ([]domain.Operator, error)
	FindByID(id uint) (*domain.Operator, error)
	Create(operator *domain.Operator) error
	Update(operator *domain.Operator) error
	UpdateStatus(id uint, status string) error
	Delete(id uint) error
}

type operatorRepository struct {
	db *gorm.DB
}

func NewOperatorRepository(db *gorm.DB) OperatorRepository {
	return &operatorRepository{db: db}
}

func (r *operatorRepository) FindAll() ([]domain.Operator, error) {
	var operators []domain.Operator
	err := r.db.Order("name asc").Find(&operators).Error
	return operators, err
}

func (r *operatorRepository) FindByID(id uint) (*domain.Operator, error) {
	var operator domain.Operator
	err := r.db.First(&operator, id).Error
	return &operator, err
}

func (r *operatorRepository) Create(operator *domain.Operator) error {
	return r.db.Create(operator).Error
}

func (r *operatorRepository) Update(operator *domain.Operator) error {
	return r.db.Save(operator).Error
}

func (r *operatorRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Operator{}).Where("id = ?", id).Update("status", status).Error
}

func (r *operatorRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Operator{}, id).Error
}

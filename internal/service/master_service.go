package service

import (
	"fmt"
	"time"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/repository"
)

type MasterService interface {
	// Customers
	GetAllCustomers() ([]domain.Customer, error)
	GetCustomerByID(id uint) (*domain.Customer, error)
	CreateCustomer(customer *domain.Customer) error
	UpdateCustomer(customer *domain.Customer) error
	DeleteCustomer(id uint) error

	// Forklift Units
	GetAllUnits(statusFilter string) ([]domain.ForkliftUnit, error)
	GetUnitByID(id uint) (*domain.ForkliftUnit, error)
	CreateUnit(unit *domain.ForkliftUnit) error
	UpdateUnit(unit *domain.ForkliftUnit) error
	DeleteUnit(id uint) error

	// Operators
	GetAllOperators() ([]domain.Operator, error)
	GetOperatorByID(id uint) (*domain.Operator, error)
	CreateOperator(operator *domain.Operator) error
	UpdateOperator(operator *domain.Operator) error
	DeleteOperator(id uint) error
}

type masterService struct {
	customerRepo repository.CustomerRepository
	unitRepo     repository.UnitRepository
	operatorRepo repository.OperatorRepository
}

func NewMasterService(
	cRepo repository.CustomerRepository,
	uRepo repository.UnitRepository,
	oRepo repository.OperatorRepository,
) MasterService {
	return &masterService{
		customerRepo: cRepo,
		unitRepo:     uRepo,
		operatorRepo: oRepo,
	}
}

// Customers
func (s *masterService) GetAllCustomers() ([]domain.Customer, error) {
	return s.customerRepo.FindAll()
}

func (s *masterService) GetCustomerByID(id uint) (*domain.Customer, error) {
	return s.customerRepo.FindByID(id)
}

func (s *masterService) CreateCustomer(customer *domain.Customer) error {
	if customer.CustomerCode == "" {
		count, _ := s.customerRepo.Count()
		customer.CustomerCode = fmt.Sprintf("CUST-%03d", count+1)
	}
	return s.customerRepo.Create(customer)
}

func (s *masterService) UpdateCustomer(customer *domain.Customer) error {
	return s.customerRepo.Update(customer)
}

func (s *masterService) DeleteCustomer(id uint) error {
	return s.customerRepo.Delete(id)
}

// Units
func (s *masterService) GetAllUnits(statusFilter string) ([]domain.ForkliftUnit, error) {
	return s.unitRepo.FindAll(statusFilter)
}

func (s *masterService) GetUnitByID(id uint) (*domain.ForkliftUnit, error) {
	return s.unitRepo.FindByID(id)
}

func (s *masterService) CreateUnit(unit *domain.ForkliftUnit) error {
	if unit.UnitCode == "" {
		unit.UnitCode = fmt.Sprintf("FL-%d", time.Now().Unix()%1000)
	}
	if unit.Status == "" {
		unit.Status = "AVAILABLE"
	}
	return s.unitRepo.Create(unit)
}

func (s *masterService) UpdateUnit(unit *domain.ForkliftUnit) error {
	return s.unitRepo.Update(unit)
}

func (s *masterService) DeleteUnit(id uint) error {
	return s.unitRepo.Delete(id)
}

// Operators
func (s *masterService) GetAllOperators() ([]domain.Operator, error) {
	return s.operatorRepo.FindAll()
}

func (s *masterService) GetOperatorByID(id uint) (*domain.Operator, error) {
	return s.operatorRepo.FindByID(id)
}

func (s *masterService) CreateOperator(operator *domain.Operator) error {
	if operator.Nip == "" {
		operator.Nip = fmt.Sprintf("OPR-%d", time.Now().Unix()%1000)
	}
	if operator.Status == "" {
		operator.Status = "STANDBY"
	}
	return s.operatorRepo.Create(operator)
}

func (s *masterService) UpdateOperator(operator *domain.Operator) error {
	return s.operatorRepo.Update(operator)
}

func (s *masterService) DeleteOperator(id uint) error {
	return s.operatorRepo.Delete(id)
}

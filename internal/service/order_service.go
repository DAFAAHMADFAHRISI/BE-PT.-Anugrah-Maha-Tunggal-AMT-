package service

import (
	"fmt"
	"time"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/repository"
)

type OrderService interface {
	GetAllOrders(status string) ([]domain.RentalOrder, error)
	GetOrderByID(id uint) (*domain.RentalOrder, error)
	CreateOrder(order *domain.RentalOrder) error
	UpdateOrderStatus(id uint, status string) error
	UpdateOrder(order *domain.RentalOrder) error
}

type orderService struct {
	orderRepo repository.RentalOrderRepository
	unitRepo  repository.UnitRepository
}

func NewOrderService(orderRepo repository.RentalOrderRepository, unitRepo repository.UnitRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		unitRepo:  unitRepo,
	}
}

func (s *orderService) GetAllOrders(status string) ([]domain.RentalOrder, error) {
	return s.orderRepo.FindAll(status)
}

func (s *orderService) GetOrderByID(id uint) (*domain.RentalOrder, error) {
	return s.orderRepo.FindByID(id)
}

func (s *orderService) CreateOrder(order *domain.RentalOrder) error {
	// Auto generate Order Number jika kosong: ORD-YYYYMM-XXXX
	if order.OrderNumber == "" {
		now := time.Now()
		count, _ := s.orderRepo.Count()
		// Coba nomor berturut-turut untuk menghindari duplikat
		for i := count + 1; i <= count+100; i++ {
			candidate := fmt.Sprintf("ORD-%s-%04d", now.Format("200601"), i)
			_, err := s.orderRepo.FindByOrderNumber(candidate)
			if err != nil {
				// Tidak ditemukan = nomor tersedia
				order.OrderNumber = candidate
				break
			}
		}
		if order.OrderNumber == "" {
			order.OrderNumber = fmt.Sprintf("ORD-%s-%04d", now.Format("200601"), count+1)
		}
	}
	if order.OrderDate == "" {
		order.OrderDate = time.Now().Format("2006-01-02")
	}
	if order.Status == "" {
		order.Status = "PENDING"
	}

	return s.orderRepo.Create(order)
}

func (s *orderService) UpdateOrderStatus(id uint, status string) error {
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *orderService) UpdateOrder(order *domain.RentalOrder) error {
	return s.orderRepo.Update(order)
}

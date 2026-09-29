package service

import (
	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/repository"
)

type DashboardMetrics struct {
	UnitStats struct {
		Total       int64 `json:"total"`
		Available   int64 `json:"available"`
		Rented      int64 `json:"rented"`
		Maintenance int64 `json:"maintenance"`
	} `json:"unit_stats"`
	OrderStats struct {
		TotalPending   int64 `json:"total_pending"`
		TotalActive    int64 `json:"total_active"`
		TotalCompleted int64 `json:"total_completed"`
	} `json:"order_stats"`
	TotalCustomers   int64                   `json:"total_customers"`
	ActiveOperations []domain.DeliveryLetter `json:"active_operations"`
	RecentOrders     []domain.RentalOrder    `json:"recent_orders"`
}

type DashboardService interface {
	GetDirekturDashboardData() (*DashboardMetrics, error)
}

type dashboardService struct {
	unitRepo     repository.UnitRepository
	orderRepo    repository.RentalOrderRepository
	customerRepo repository.CustomerRepository
	deliveryRepo repository.DeliveryLetterRepository
}

func NewDashboardService(
	uRepo repository.UnitRepository,
	oRepo repository.RentalOrderRepository,
	cRepo repository.CustomerRepository,
	dRepo repository.DeliveryLetterRepository,
) DashboardService {
	return &dashboardService{
		unitRepo:     uRepo,
		orderRepo:    oRepo,
		customerRepo: cRepo,
		deliveryRepo: dRepo,
	}
}

func (s *dashboardService) GetDirekturDashboardData() (*DashboardMetrics, error) {
	metrics := &DashboardMetrics{}

	// Statistik Unit
	tot, avail, rent, maint, _ := s.unitRepo.GetStatistics()
	metrics.UnitStats.Total = tot
	metrics.UnitStats.Available = avail
	metrics.UnitStats.Rented = rent
	metrics.UnitStats.Maintenance = maint

	// Statistik Customer
	custCount, _ := s.customerRepo.Count()
	metrics.TotalCustomers = custCount

	// Statistik Order
	allOrders, _ := s.orderRepo.FindAll("")
	for _, o := range allOrders {
		switch o.Status {
		case "PENDING":
			metrics.OrderStats.TotalPending++
		case "ACTIVE", "APPROVED":
			metrics.OrderStats.TotalActive++
		case "COMPLETED":
			metrics.OrderStats.TotalCompleted++
		}
	}

	// 5 Recent Orders
	if len(allOrders) > 5 {
		metrics.RecentOrders = allOrders[:5]
	} else {
		metrics.RecentOrders = allOrders
	}

	// Active Operations (ASSIGNED, ON_THE_WAY, WORKING)
	allLetters, _ := s.deliveryRepo.FindAll("")
	var activeOps []domain.DeliveryLetter
	for _, l := range allLetters {
		if l.OperationalStatus == "WORKING" || l.OperationalStatus == "ON_THE_WAY" || l.OperationalStatus == "ASSIGNED" {
			activeOps = append(activeOps, l)
		}
	}
	metrics.ActiveOperations = activeOps

	return metrics, nil
}

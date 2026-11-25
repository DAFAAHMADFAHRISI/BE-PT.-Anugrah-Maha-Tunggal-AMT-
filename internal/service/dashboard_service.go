package service

import (
	"fmt"
	"sort"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/repository"
)

type PeriodYear struct {
	Year   int   `json:"year"`
	Months []int `json:"months"`
}

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
	AllOrders        []domain.RentalOrder    `json:"all_orders"`
	AllOperations    []domain.DeliveryLetter `json:"all_operations"`
	AvailablePeriods []PeriodYear            `json:"available_periods"`
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

func extractAvailablePeriods(orders []domain.RentalOrder, letters []domain.DeliveryLetter) []PeriodYear {
	yearMap := make(map[int]map[int]bool)

	addDate := func(dateStr string) {
		if len(dateStr) >= 7 {
			var y, m int
			if n, _ := fmt.Sscanf(dateStr[:7], "%d-%d", &y, &m); n == 2 && y > 2000 {
				if _, ok := yearMap[y]; !ok {
					yearMap[y] = make(map[int]bool)
				}
				yearMap[y][m] = true
			}
		}
	}

	for _, o := range orders {
		if o.StartDate != "" {
			addDate(o.StartDate)
		} else if o.OrderDate != "" {
			addDate(o.OrderDate)
		}
	}

	for _, l := range letters {
		if l.IssueDate != "" {
			addDate(l.IssueDate)
		}
	}

	// Urutkan tahun descending (misal 2026, 2022)
	var years []int
	for y := range yearMap {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	var result []PeriodYear
	for _, y := range years {
		var months []int
		for m := range yearMap[y] {
			months = append(months, m)
		}
		sort.Ints(months)
		result = append(result, PeriodYear{
			Year:   y,
			Months: months,
		})
	}

	return result
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

	// Statistik Order & Seluruh Data Order
	allOrders, _ := s.orderRepo.FindAll("")
	metrics.AllOrders = allOrders
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

	// 10 Recent Orders untuk list ringkasan
	if len(allOrders) > 10 {
		metrics.RecentOrders = allOrders[:10]
	} else {
		metrics.RecentOrders = allOrders
	}

	// Seluruh Surat Jalan Operasional
	allLetters, _ := s.deliveryRepo.FindAll("")
	metrics.AllOperations = allLetters
	var activeOps []domain.DeliveryLetter
	for _, l := range allLetters {
		if l.OperationalStatus == "WORKING" || l.OperationalStatus == "ON_THE_WAY" || l.OperationalStatus == "ASSIGNED" {
			activeOps = append(activeOps, l)
		}
	}
	metrics.ActiveOperations = activeOps

	// Ekstrak periode tahun & bulan yang riil ada di database
	metrics.AvailablePeriods = extractAvailablePeriods(allOrders, allLetters)

	return metrics, nil
}

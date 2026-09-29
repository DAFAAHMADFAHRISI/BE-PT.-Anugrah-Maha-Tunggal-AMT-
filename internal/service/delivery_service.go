package service

import (
	"errors"
	"fmt"
	"time"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/repository"
)

type DeliveryService interface {
	GetAllDeliveryLetters(status string) ([]domain.DeliveryLetter, error)
	GetDeliveryLetterByID(id uint) (*domain.DeliveryLetter, error)
	CreateDeliveryLetter(letter *domain.DeliveryLetter) error
	UpdateOperationalStatus(id uint, newStatus, notes string) error
	GetMovementLogs(unitID uint) ([]domain.UnitMovementLog, error)
}

type deliveryService struct {
	deliveryRepo repository.DeliveryLetterRepository
	orderRepo    repository.RentalOrderRepository
	unitRepo     repository.UnitRepository
	operatorRepo repository.OperatorRepository
}

func NewDeliveryService(
	dRepo repository.DeliveryLetterRepository,
	oRepo repository.RentalOrderRepository,
	uRepo repository.UnitRepository,
	opRepo repository.OperatorRepository,
) DeliveryService {
	return &deliveryService{
		deliveryRepo: dRepo,
		orderRepo:    oRepo,
		unitRepo:     uRepo,
		operatorRepo: opRepo,
	}
}

func (s *deliveryService) GetAllDeliveryLetters(status string) ([]domain.DeliveryLetter, error) {
	return s.deliveryRepo.FindAll(status)
}

func (s *deliveryService) GetDeliveryLetterByID(id uint) (*domain.DeliveryLetter, error) {
	return s.deliveryRepo.FindByID(id)
}

func (s *deliveryService) CreateDeliveryLetter(letter *domain.DeliveryLetter) error {
	// Validasi order
	order, err := s.orderRepo.FindByID(letter.RentalOrderID)
	if err != nil {
		return errors.New("pesanan sewa (rental order) tidak ditemukan")
	}

	// Validasi ketersediaan unit
	unit, err := s.unitRepo.FindByID(letter.UnitID)
	if err != nil {
		return errors.New("unit forklift tidak ditemukan")
	}
	if unit.Status == "MAINTENANCE" {
		return errors.New("unit forklift sedang dalam perbaikan (MAINTENANCE) dan tidak dapat dialokasikan")
	}

	// Auto generate Nomor Surat Jalan jika kosong: SJ-YYYYMM-XXXX
	if letter.LetterNumber == "" {
		now := time.Now()
		letters, _ := s.deliveryRepo.FindAll("")
		letter.LetterNumber = fmt.Sprintf("SJ-%s-%04d", now.Format("200601"), len(letters)+1)
	}

	if letter.IssueDate == "" {
		letter.IssueDate = time.Now().Format("2006-01-02")
	}

	if letter.OperationalStatus == "" {
		letter.OperationalStatus = "ASSIGNED"
	}

	// Simpan Surat Jalan
	if err := s.deliveryRepo.Create(letter); err != nil {
		return err
	}

	// Update status Unit menjadi RENTED
	_ = s.unitRepo.UpdateStatus(letter.UnitID, "RENTED")

	// Update status Operator menjadi ON_DUTY
	_ = s.operatorRepo.UpdateStatus(letter.OperatorID, "ON_DUTY")

	// Update status Order menjadi ACTIVE
	_ = s.orderRepo.UpdateStatus(order.ID, "ACTIVE")

	// Catat log pergerakan awal
	_ = s.deliveryRepo.CreateMovementLog(&domain.UnitMovementLog{
		UnitID:           letter.UnitID,
		DeliveryLetterID: &letter.ID,
		StatusFrom:       "AVAILABLE",
		StatusTo:         "ASSIGNED",
		Notes:            fmt.Sprintf("Surat Jalan %s diterbitkan untuk pesanan %s", letter.LetterNumber, order.OrderNumber),
		LoggedAt:         time.Now(),
	})

	return nil
}

func (s *deliveryService) UpdateOperationalStatus(id uint, newStatus, notes string) error {
	letter, err := s.deliveryRepo.FindByID(id)
	if err != nil {
		return errors.New("surat jalan tidak ditemukan")
	}

	oldStatus := letter.OperationalStatus
	if err := s.deliveryRepo.UpdateOperationalStatus(id, newStatus); err != nil {
		return err
	}

	// Log audit trail pergerakan
	logEntry := domain.UnitMovementLog{
		UnitID:           letter.UnitID,
		DeliveryLetterID: &letter.ID,
		StatusFrom:       oldStatus,
		StatusTo:         newStatus,
		Notes:            notes,
		LoggedAt:         time.Now(),
	}
	_ = s.deliveryRepo.CreateMovementLog(&logEntry)

	// Jika status berubah menjadi FINISHED (Pekerjaan Selesai)
	if newStatus == "FINISHED" {
		// Unit kembali AVAILABLE di pool
		_ = s.unitRepo.UpdateStatus(letter.UnitID, "AVAILABLE")
		// Operator kembali STANDBY
		_ = s.operatorRepo.UpdateStatus(letter.OperatorID, "STANDBY")
		// Order menjadi COMPLETED
		_ = s.orderRepo.UpdateStatus(letter.RentalOrderID, "COMPLETED")

		_ = s.deliveryRepo.CreateMovementLog(&domain.UnitMovementLog{
			UnitID:           letter.UnitID,
			DeliveryLetterID: &letter.ID,
			StatusFrom:       "RENTED",
			StatusTo:         "AVAILABLE",
			Notes:            "Operasional selesai, unit kembali ke pool dan siap disewa kembali",
			LoggedAt:         time.Now(),
		})
	}

	return nil
}

func (s *deliveryService) GetMovementLogs(unitID uint) ([]domain.UnitMovementLog, error) {
	return s.deliveryRepo.GetMovementLogs(unitID)
}

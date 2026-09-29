package service

import (
	"errors"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/repository"
	"erp-anugrah-maha-tunggal/pkg/utils"
)

type AuthService interface {
	Login(username, password string) (string, *domain.User, error)
	Register(user *domain.User, plainPassword string) error
	GetProfile(id uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.FindAll()
}

func (s *authService) Login(username, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	if !user.IsActive {
		return "", nil, errors.New("akun Anda dinonaktifkan, silakan hubungi administrator")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("username atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), user.Name)
	if err != nil {
		return "", nil, errors.New("gagal menerbitkan token otorisasi")
	}

	return token, user, nil
}

func (s *authService) Register(user *domain.User, plainPassword string) error {
	hashed, err := utils.HashPassword(plainPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hashed
	return s.userRepo.Create(user)
}

func (s *authService) GetProfile(id uint) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

package repository

import (
	"erp-anugrah-maha-tunggal/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	FindAll() ([]domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ? OR email = ?", username, username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}

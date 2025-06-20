package user

import (
	"ecommerce-stock-api/models"
	"ecommerce-stock-api/pkg/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(email, phone, password string) error
	Login(identifier, password string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Register(email, phone, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user := &models.User{
		Email:        email,
		Phone:        phone,
		PasswordHash: hashedPassword,
	}
	err = s.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Login(identifier, password string) (string, error) {
	user, err := s.repo.FindByEmailOrPhone(identifier)
	if err != nil {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

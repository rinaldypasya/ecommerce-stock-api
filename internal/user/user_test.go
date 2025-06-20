package user_test

import (
	"ecommerce-stock-api/internal/user"
	"ecommerce-stock-api/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	users     []models.User
	findErr   bool
	createErr bool
}

func (m *mockRepo) FindByEmailOrPhone(identifier string) (*models.User, error) {
	if m.findErr {
		return nil, errors.New("not found")
	}
	for _, u := range m.users {
		if u.Email == identifier || u.Phone == identifier {
			return &u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) Create(user *models.User) error {
	if m.createErr {
		return errors.New("create error")
	}
	m.users = append(m.users, *user)
	return nil
}

func TestRegisterSuccess(t *testing.T) {
	mock := &mockRepo{}
	svc := user.NewService(mock)

	err := svc.Register("test@example.com", "1234567890", "secret")
	assert.NoError(t, err)
	assert.Len(t, mock.users, 1)
	assert.Equal(t, "test@example.com", mock.users[0].Email)
}

func TestRegisterFailure(t *testing.T) {
	mock := &mockRepo{createErr: true}
	svc := user.NewService(mock)

	err := svc.Register("fail@example.com", "000", "secret")
	assert.Error(t, err)
}

func TestLoginSuccess(t *testing.T) {
	hashed, _ := user.HashPassword("secret")
	mock := &mockRepo{
		users: []models.User{
			{Email: "me@example.com", PasswordHash: hashed},
		},
	}
	svc := user.NewService(mock)

	token, err := svc.Login("me@example.com", "secret")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginWrongPassword(t *testing.T) {
	hashed, _ := user.HashPassword("secret")
	mock := &mockRepo{
		users: []models.User{
			{Email: "me@example.com", PasswordHash: hashed},
		},
	}
	svc := user.NewService(mock)

	_, err := svc.Login("me@example.com", "wrongpass")
	assert.Error(t, err)
}

func TestLoginUserNotFound(t *testing.T) {
	mock := &mockRepo{}
	svc := user.NewService(mock)

	_, err := svc.Login("unknown@example.com", "any")
	assert.Error(t, err)
}

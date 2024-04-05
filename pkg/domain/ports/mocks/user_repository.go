package mocks

import (
	"snippetbox/pkg/domain/models"
	"time"
)

var mockUser = &models.User{
	ID:      1,
	Name:    "Alice",
	Email:   "alice@example.com",
	Created: time.Now(),
	Active:  true,
}

type MockUserRepository struct{}

func (u MockUserRepository) Create(e models.User) (uint, error) {

	switch e.Email {
	case "dupe@example.com":
		return 0, models.ErrDuplicateEmail
	default:
		return uint(e.ID), nil
	}

}

func (u MockUserRepository) ReadAll() ([]*models.User, error) {
	return []*models.User{mockUser}, nil
}

func (u MockUserRepository) ReadOne(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (u MockUserRepository) ReadBy(e *models.User) ([]*models.User, error) {
	if e != nil {
		return []*models.User{mockUser}, nil
	}
	return nil, models.ErrNoRecord
}

func (u MockUserRepository) Update(e *models.User) (uint, error) {
	if e == nil {
		return 0, models.ErrNoRecord
	}
	mockUser.ID = e.ID
	mockUser.Name = e.Name
	mockUser.Email = e.Email
	mockUser.Created = e.Created
	mockUser.Active = e.Active
	return uint(mockUser.ID), nil
}

func (u MockUserRepository) Delete(id uint) (uint, error) {
	switch id {
	case 1:
		return id, nil
	default:
		return 0, models.ErrNoRecord
	}
}

func (u MockUserRepository) Authenticate(email, password string) (int, error) {
	if email == "" || password == "" {
		return 0, models.ErrInvalidCredentials
	} else if email == "mail@example.com" && password == "123" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

package types

import (
	"github.com/Ayobami6/pickitup_v3/pkg/models"
)

type UserRepo interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

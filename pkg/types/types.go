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

type RiderRepo interface {
	CreateRider(rider *models.Rider) error
	GetRiders() (*[]models.Rider, error)
	GetRiderByID(id uint) (*models.Rider, error)
	GetRiderByUserID(userID uint) (*models.Rider, error)
	GetRiderReviews(riderID uint) (*[]models.Review, error)
	UpdateRiderRating(riderID uint) error
	UpdateRiderMinAndMaxCharge(minCharge float64, maxCharge float64, userID uint) error
	UpdateRiderAvailability(riderID uint, status models.RiderAvailabilityStatus) error
	Save(rider *models.Rider) error
}


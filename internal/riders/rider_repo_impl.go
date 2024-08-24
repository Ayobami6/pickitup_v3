package riders

import (
	"math"

	"github.com/Ayobami6/pickitup_v3/pkg/models"
	"gorm.io/gorm"
)

type RideRepoImpl struct {
	db *gorm.DB
}


func NewRideRepoImpl(db *gorm.DB) *RideRepoImpl {
    return &RideRepoImpl{db: db}
}

func (r *RideRepoImpl) CreateRide(rider *models.Rider) error {
    return r.db.Create(rider).Error
}

func (r *RideRepoImpl)GetRiders()(*[]models.Rider, error) {
	riders := []models.Rider{}
    res := r.db.Find(&riders)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &riders, nil
}

func (r *RideRepoImpl) GetRideById(id uint) (*models.Rider, error) {
	// get rider by id
	var rider models.Rider
    res := r.db.First(&rider, id)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &rider, nil
}

func (r *RideRepoImpl)GetRiderByUserID(userID uint) (*models.Rider, error) {
	// get rider by user id
    var rider models.Rider
    res := r.db.Where(&models.Rider{UserID: uint(userID)}).First(&rider)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &rider, nil
}

func (r *RideRepoImpl)GetRiderReviews(riderId uint)(*[]models.Review,error) {
	// get all reviews for the rider
    var reviews []models.Review
    res := r.db.Where(&models.Review{RiderID: uint(riderId)}).Find(&reviews)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &reviews, nil
}

func (r *RideRepoImpl)UpdateRiderRating(riderID uint) error {
	// get rider by id
	rider, err := r.GetRideById(riderID)
	if err!= nil {
        return err
    }
	// get rider reviews
	reviews, err := r.GetRiderReviews(riderID)
    if err!= nil {
        return err
    }
    // calculate average rating
    var totalRating float64 = 0
    for _, review := range *reviews {
        totalRating += review.Rating
    }
    averageRating := totalRating / float64(len(*reviews))
	// round average rating
	averageRating = (math.Round(averageRating*10) / 10)
    rider.Rating = averageRating
    return r.db.Save(rider).Error

}

func (r *RideRepoImpl)Save(rider models.Rider) error {
	return r.db.Save(&rider).Error
}

func (r * RideRepoImpl)UpdateMinAndMaxCharge(minCharge float64, maxCharge float64, userID uint)error {
	res := r.db.Where(&models.Rider{UserID: uint(userID)}).Updates(&models.Rider{MinimumCharge: minCharge, MaximumCharge: maxCharge})
    if res.Error!= nil {
        return res.Error
    }
    return nil
}

func (r *RideRepoImpl) UpdateRiderAvailability(riderId uint, status models.RiderAvailabilityStatus) error {
	res := r.db.Where(&models.Rider{UserID: uint(riderId)}).Updates(&models.Rider{AvailabilityStatus: status})
    if res.Error!= nil {
        return res.Error
    }
    return nil
}
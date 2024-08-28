package riders

import (
	"errors"
	"fmt"
	"log"

	"github.com/Ayobami6/pickitup_v3/internal/riders/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/auth"
	"github.com/Ayobami6/pickitup_v3/pkg/models"
	"github.com/Ayobami6/pickitup_v3/pkg/types"
	"github.com/Ayobami6/pickitup_v3/pkg/utils"
)

type RiderService struct {
	riderRepo types.RiderRepo
	userRepo types.UserRepo
}

func NewRiderService(repo types.RiderRepo, userRepo types.UserRepo) *RiderService {
    return &RiderService{riderRepo: repo, userRepo: userRepo}
}

func (rs *RiderService)CreateRider(pl *dto.RegisterRiderDTO) error {
	// extract all data from payload
	email := pl.Email
	password := pl.Password
	phoneNumber := pl.PhoneNumber
	address := pl.Address
	bikeNumber := pl.BikeNumber
	driverLicenseNumber := pl.DriverLicenseNumber
	nextOfKinName := pl.NextOfKinName
	nextOfKinPhone := pl.NextOfKinPhone
	nextOfKinAddress := pl.NextOfKinAddress
	userName := pl.UserName
	firstName := pl.FirstName
	lastName := pl.LastName
	// check if user exists already
	_, err := rs.userRepo.GetUserByEmail(email)
	if err == nil {
        return errors.New("user with this email already exists")
    }
	// hash password
	hashedPassword, err := auth.HashPassword(password)
	if err!= nil {
		return utils.ThrowError(err)
    }
	// create user with the new hashed password
	user := models.User{
		Email:        email,
        Password:     hashedPassword,
        UserName:     userName,
        PhoneNumber: phoneNumber,
	}
	if err := rs.userRepo.CreateUser(&user); err!= nil {
		log.Printf("Error creating user %v \n", err.Error())
		return utils.ThrowError(err)
    }
	// create rider with the extracted data
	rider := models.Rider{
        UserID:             user.ID,
        FirstName:          firstName,
        LastName:           lastName,
        Address:            address,
        BikeNumber:         bikeNumber,
        DriverLicenseNumber: driverLicenseNumber,
        NextOfKinName:      nextOfKinName,
        NextOfKinPhone:     nextOfKinPhone,
        NextOfKinAddress:   nextOfKinAddress,
    }
    if err := rs.riderRepo.CreateRider(&rider); err!= nil {
        return utils.ThrowError(err)
    }
	// generate verification number
	num, err := utils.GenerateAndCacheVerificationCode(email)
	if err!= nil {
        log.Println("Generate Code failed: ", err)
    } else {
		// send the email to verify
		msg := fmt.Sprintf("Your verification code is %d\n", num)
        go utils.SendMail(email, "Email Verification", userName, msg)
	}
    return nil

}

func (rs *RiderService)GetRiders() (*[]dto.RiderListResponse, error) {
	// get riders
	riders, err := rs.riderRepo.GetRiders();
	var riderDtoList []dto.RiderListResponse
	if err!= nil {
        log.Printf("Couldn't fetch riders %v \n", err)
        return nil, utils.ThrowError(err)
    }
	riderList := *riders
	for _, rider := range riderList {
		riderDto := dto.RiderListResponse{
            RiderID:             rider.RiderID,
            FirstName:          rider.FirstName,
            LastName:           rider.LastName,
            Address:            rider.Address,
            BikeNumber:         rider.BikeNumber,
            Rating: rider.Rating,
            SuccessfulRides:     rider.SuccessfulRides,
            Level:              rider.Level,
            CurrentLocation:    rider.CurrentLocation,
            AvailabilityStatus:     string(rider.AvailabilityStatus),
            MaximumCharge:      rider.MaximumCharge,
            MinimumCharge:      rider.MinimumCharge,
			ID: rider.ID,
        }
        riderDtoList = append(riderDtoList, riderDto)
    }
	return &riderDtoList, nil

}

func (rs *RiderService)GetRider(riderID uint) (*dto.RiderResponse, error) {
	// get rider by id
    rider, err := rs.riderRepo.GetRiderByID(riderID);
    if err!= nil {
        log.Printf("Couldn't fetch rider %v \n", err)
        return nil, utils.ThrowError(err)
    }
	// get reviews 
	var reviewsDto []dto.ReviewResponse
	reviews, err := rs.riderRepo.GetRiderReviews(riderID)
	if err!= nil {
        log.Printf("Couldn't fetch reviews for rider %v \n", err)
        
    } else {

		for _, review := range *reviews {
			reviewsDto = append(reviewsDto, dto.ReviewResponse{
				Rating: review.Rating,
				Comment: review.Comment,
			})
		}
	}
	res := dto.RiderResponse{
		ID: rider.ID,
		RiderID:             rider.RiderID,
        FirstName:          rider.FirstName,
        LastName:           rider.LastName,
        Address:            rider.Address,
        BikeNumber:         rider.BikeNumber,
		Rating: rider.Rating,
        SuccessfulRides:     rider.SuccessfulRides,
		Level:              rider.Level,
        CurrentLocation:    rider.CurrentLocation,
        AvailabilityStatus:     string(rider.AvailabilityStatus),
		MaximumCharge:      rider.MaximumCharge,
        MinimumCharge:      rider.MinimumCharge,
		Reviews: reviewsDto,
	}
    return &res, nil

}

func (rs *RiderService)UpdateCharges(pl *dto.UpdateChargeDTO, userId uint) error {
	minCharge := pl.MinimumCharge
    maxCharge := pl.MaximumCharge
    userID := userId
    // update charges in the database
    err := rs.riderRepo.UpdateRiderMinAndMaxCharge(minCharge, maxCharge, userID)
    if err!= nil {
        return utils.ThrowError(err)
    }
    return nil
}

func (rs *RiderService)UpdateRiderAvailability(pl *dto.UpdateRiderAvailabilityStatusDTO, userId uint) error {
	availabilityStatus := pl.AvailabilityStatus
    userID := userId
    // update availability status in the database
    err := rs.riderRepo.UpdateRiderAvailability(userID, models.RiderAvailabilityStatus(availabilityStatus))
    if err!= nil {
        return utils.ThrowError(err)
    }
    return nil
}
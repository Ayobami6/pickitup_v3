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

func (rs *RiderService)CreateRider(pl dto.RegisterRiderDTO) error {
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
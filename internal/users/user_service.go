package users

import (
	"errors"
	"fmt"
	"log"

	"github.com/Ayobami6/pickitup_v3/config"
	"github.com/Ayobami6/pickitup_v3/internal/users/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/auth"
	"github.com/Ayobami6/pickitup_v3/pkg/models"
	"github.com/Ayobami6/pickitup_v3/pkg/types"
	"github.com/Ayobami6/pickitup_v3/pkg/utils"
)

type UserService struct {
	repo types.UserRepo
}

// constructor
func NewUserService(repo types.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService)RegisterUser(pl dto.RegisterUserDTO) (any, error) {
	// validate input
    email := pl.Email
	password := pl.Password
	username := pl.UserName
	phone_number := pl.PhoneNumber
	// check if user is already registered
	_, err := u.repo.GetUserByEmail(email)
    if err == nil {
        return nil, errors.New("user with this email already exists")
    }
    // hash password
    hashedPassword, err := auth.HashPassword(password)
    if err!= nil {
        return "", err
    }
    // create user
    user := &models.User{Email: email, Password: hashedPassword, UserName: username, PhoneNumber: phone_number}
    err = u.repo.CreateUser(user)
    if err!= nil {
        return "", err
    }
	// implement send otp for verification
	num, err := utils.GenerateAndCacheVerificationCode(email)
	if err != nil {
		log.Println("Generate Code Failed: ", err)
	} else {
		// send the mail
		msg := fmt.Sprintf("Your verification code is %d\n", num)
		go utils.SendMail(email, "Email Verification", username, msg)
	}
	message := "Registration Successfully"
    return message, nil
}


func (u *UserService) LoginUser(pl dto.LoginDTO) (map[string]string, error) {
	email := pl.Email
    password := pl.Password
    // check if user exists
    user, err := u.repo.GetUserByEmail(email)
    if err!= nil {
        return nil, errors.New("user not found")
    }
    // compare hashed passwords
    if !auth.CheckPassword(user.Password, []byte(password)) {
        return nil, errors.New("invalid password")
    }
    // generate and return JWT token
	secret := []byte(config.GetEnv("JWT_SECRET", "secret"))
    token, err := auth.CreateJWT(secret, int(user.ID))
    if err!= nil {
		log.Println("Failed to create jwt token: ", err)
        return nil, err
    }
    return token, nil
}
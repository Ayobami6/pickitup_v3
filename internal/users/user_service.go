package users

import (
	"errors"
	"fmt"

	"github.com/Ayobami6/pickitup_v3/internal/users/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/auth"
	"github.com/Ayobami6/pickitup_v3/pkg/models"
	"github.com/Ayobami6/pickitup_v3/pkg/types"
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
	message := fmt.Sprintf("Registration Successfully")
    return message, nil
}
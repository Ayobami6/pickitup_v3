package dto

type RegisterUserDTO struct {
	UserName    string `json:"username" binding:"required"`
	Password    string `json:"password" validate:"required,min=6, binding:"required,min=6,alphanumunicode"`
	Email       string `json:"email" validate:"required,email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required" binding:"required"`
}
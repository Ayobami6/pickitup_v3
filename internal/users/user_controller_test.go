package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ayobami6/pickitup_v3/internal/users/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestUserRegister(t *testing.T) {
	userRepo:= &mockUserRepo{}
	userService := NewUserService(userRepo)
	userController := NewUserController(*userService)

	t.Run("Should fail for bad data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		payload := dto.RegisterUserDTO{}
		jsonPayload, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		userController.RegisterUser(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should fail if password is less than 6", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
        payload := dto.RegisterUserDTO{
			Email: "ayobamidele@gmail.com", 
			UserName: "first_name", PhoneNumber: "070224754332", Password: "123"}
        jsonPayload, _ := json.Marshal(payload)
        req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonPayload))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = req
	})


}


// lets mock userService 

type mockUserRepo struct {}

// implement all interface methods

func (m *mockUserRepo) CreateUser(user *models.User) (error){
	return nil
}

func (u *mockUserRepo) GetUserByEmail(email string) (*models.User, error) {
    return nil, nil
}

func (u *mockUserRepo)GetUserByID(id uint) (*models.User, error) {
    return nil, nil
}

func (u *mockUserRepo)UpdateUser( updatedUser *models.User) error {
    // update a user
    return nil
}
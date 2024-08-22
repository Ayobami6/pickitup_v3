package users

import (
	"log"

	"github.com/Ayobami6/pickitup_v3/pkg/models"
	"gorm.io/gorm"
)

// object structure aka class
type UserRepoImpl struct {
	db *gorm.DB
}

// class constructor for DI
func NewUserRepoImpl(db *gorm.DB) *UserRepoImpl {
	// migrate table
	err := db.AutoMigrate(&models.User{})
    if err!= nil {
        log.Fatal(err)
    }
    return &UserRepoImpl{db: db}
}

// methods implementations

func (u *UserRepoImpl) CreateUser(user *models.User) (error) {
	res := u.db.Create(&user)
	if res.Error!= nil {
        return res.Error
    }
	return nil
}


func (u *UserRepoImpl) GetUserByEmail(email string) (*models.User, error) {
	result := &models.User{}
    err := u.db.Where("email =?", email).First(&result).Error
    if err!= nil {
        return nil, err
    }
    return result, nil
}

func (u *UserRepoImpl) GetUserByID(id uint) (*models.User, error) {
	result := &models.User{}
    res := u.db.First(&result, id)
    if res.Error!= nil {
        return nil, res.Error
    }
    return result, nil
}

func (u *UserRepoImpl) UpdateUser( updatedUser *models.User) error {
    // update a user
    res := u.db.Save(&updatedUser)
    if res.Error!= nil {
        return res.Error
    }
    return nil
}
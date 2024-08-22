package users

import (
	"log"
	"net/http"
	"strings"

	"github.com/Ayobami6/pickitup_v3/internal/users/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{service: service}
}


func (uc *UserController)RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.POST("/register", uc.RegisterUser)
	users.POST("/login", uc.Login)
}

func (uc *UserController)RegisterUser(c *gin.Context){
	var payload dto.RegisterUserDTO
	// bind json payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println(err)
		c.JSON(400, utils.Response(400,  nil, err.Error()))
        return
	}
	// register user 
	msg, err := uc.service.RegisterUser(payload)
    if err!= nil {
		log.Println(err)
		if strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusConflict, utils.Response(http.StatusConflict,  nil, "User with this email already exists"))
            return
		} else if strings.Contains(err.Error(),"uni_users_phone_number"){
			c.JSON(http.StatusConflict, utils.Response(http.StatusConflict, nil, "User with this phone already exists"))
			return
		}
        c.JSON(500, utils.Response(500,  nil, nil))
        return
    }
    c.JSON(201, utils.Response(http.StatusCreated,  nil, msg.(string)))
}

func (uc *UserController)Login(c *gin.Context) {
	var payload dto.LoginDTO
    // bind json payload
    if err := c.ShouldBindJSON(&payload); err!= nil {
        log.Println(err)
        c.JSON(400, utils.Response(400,  nil, err.Error()))
        return
    }
    // login user
    token, err := uc.service.LoginUser(payload)
    if err!= nil {
        log.Println(err)
        c.JSON(401, utils.Response(401,  nil, "Invalid credentials"))
        return
    }
	data := map[string]any{
		"token": token,
	}
    c.JSON(200, utils.Response(http.StatusOK, data, "Login Successful"))
}
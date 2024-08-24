package riders

import (
	"log"
	"net/http"

	"github.com/Ayobami6/pickitup_v3/internal/riders/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/utils"
	"github.com/gin-gonic/gin"
)

// rider class implementation
type RiderController struct {
	// depends on riderService
	riderService RiderService
}

// RiderController constructor

func NewRiderController(riderService RiderService) *RiderController {
    return &RiderController{riderService}
}

// controller routes registory
func (c *RiderController)RegisterRoutes(router *gin.RouterGroup){
	riders := router.Group("/riders")
	riders.POST("/register", c.RegisterRider)
}

// RegisterRider handles POST request to register a new rider

func (c *RiderController) RegisterRider(ctx *gin.Context) {
    // implement the logic to register a new rider
	var pl dto.RegisterRiderDTO
	if err := ctx.ShouldBindJSON(&pl); err!= nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
	err := c.riderService.CreateRider(pl)
	if err!= nil {
		// TODO: need to handler potential errors
		log.Println(err)
        ctx.JSON(500, utils.Response(500, nil, err.Error()))
        return
    }
	ctx.JSON(201, utils.Response(http.StatusCreated, nil, "Rider Successfully Created!"))
}

// GetRider retrieves a rider by id from the service

// GetRiders retrieves all riders from the service

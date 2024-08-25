package riders

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/Ayobami6/pickitup_v3/cmd/docs"
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
	riders.GET("/:id", c.GetRider)
	riders.GET("", c.GetRiders)
}

// RegisterRider handles POST request to register a new rider

// @Summary      Register a new rider
// @Description  Register a new rider with the provided details
// @Tags         riders
// @Accept       json
// @Produce      json
// @Param        rider  body      dto.RegisterRiderDTO  true  "Rider registration details"
// @Success      201    {object}  utils.ResponseMessage
// @Failure      400    {object}  utils.ResponseMessage
// @Failure      500    {object}  utils.ResponseMessage
// @Router       /riders/register [post]
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

func (c *RiderController) GetRider(ctx *gin.Context) {
    id := ctx.Param("id")
	// convert ascii to integer
	rid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(400, utils.Response(http.StatusBadRequest, nil, err.Error()))
	}
    rider, err := c.riderService.GetRider(uint(rid))
    if err!= nil {
        ctx.JSON(404, utils.Response(http.StatusNotFound, nil, "Rider not found"))
        return
    }
	domain := utils.GetDomainUrl(ctx)
	rider.SelfUrl = fmt.Sprintf("%s/riders/%s", domain, rider.RiderID)
    ctx.JSON(200, utils.Response(http.StatusOK, rider, "Rider Fetch Successfully"))
}

// GetRiders retrieves all riders from the service

func (c *RiderController) GetRiders(ctx *gin.Context) {
    riders, err := c.riderService.GetRiders()
    if err!= nil {
        ctx.JSON(500, utils.Response(http.StatusInternalServerError, nil, err.Error()))
        return
    }
    domain := utils.GetDomainUrl(ctx)
	riderList := *riders
    for i := range riderList {
        riderList[i].SelfUrl = fmt.Sprintf("%s/riders/%s", domain, riderList[i].RiderID)
    }
    ctx.JSON(200, utils.Response(http.StatusOK, riders, "Riders Fetch Successfully"))
}

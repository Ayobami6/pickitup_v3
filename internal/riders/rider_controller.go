package riders

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/Ayobami6/pickitup_v3/cmd/docs"
	"github.com/Ayobami6/pickitup_v3/internal/riders/dto"
	"github.com/Ayobami6/pickitup_v3/pkg/auth"
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
	riders.PATCH("/charges", auth.RiderAuth(c.riderService.riderRepo), c.UpdateCharge)
	riders.PATCH("/status", auth.RiderAuth(c.riderService.riderRepo), c.UpdateStatus)
}

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
	err := c.riderService.CreateRider(&pl)
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
        riderList[i].SelfUrl = fmt.Sprintf("%s/riders/%d", domain, riderList[i].ID)
    }
    ctx.JSON(200, utils.Response(http.StatusOK, riders, "Riders Fetch Successfully"))
}


func(c *RiderController)UpdateCharge(ctx *gin.Context) {
	// get user id from context
	userId := auth.GetUserIDFromContext(ctx)
	if userId == -1 {
        auth.Forbidden(ctx)
        return
    }
	// payload
	var pl dto.UpdateChargeDTO
	if err := ctx.ShouldBindJSON(&pl); err!= nil {
        ctx.JSON(400, utils.Response(400, nil, err.Error()))
        return
    }
	// send payload to the service
	var userID uint = uint(userId)
	err := c.riderService.UpdateCharges(&pl, userID)
	if err!= nil {
        log.Println(err)
        ctx.JSON(500, utils.Response(500, nil, err.Error()))
        return
    }

	ctx.JSON(200, utils.Response(http.StatusOK, nil, "Charges Updated Successfully"))

}

func (c *RiderController)UpdateStatus(ctx *gin.Context) {
	// get user id from context
    userId := auth.GetUserIDFromContext(ctx)
    if userId == -1 {
        auth.Forbidden(ctx)
        return
    }
    // payload
    var pl *dto.UpdateRiderAvailabilityStatusDTO
	// bind json payload
	if err := ctx.ShouldBindJSON(&pl); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, nil, err.Error()))
		return
	}

	// validate status type
	var userID uint = uint(userId)

	statusMap := map[string]bool{
		"Available": true,
        "Unavailable": true,
        "On Break": true,
        "Busy": true,
	}
	if !statusMap[pl.AvailabilityStatus]{
		ctx.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, nil, "Invalid availability status"))
        return
	}
	// send payload to the service
	err := c.riderService.UpdateRiderAvailability(pl, userID)
    if err!= nil {
        log.Println(err)
        ctx.JSON(http.StatusInternalServerError, utils.Response(http.StatusInternalServerError, nil, err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.Response(http.StatusOK, nil, "Status Updated Successfully"))

}
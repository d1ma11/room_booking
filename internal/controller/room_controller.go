package controller

import (
	"net/http"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomService service.RoomService
}

func NewRoomController(roomService service.RoomService) *RoomController {
	return &RoomController{roomService: roomService}
}

type createRoomRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

func (c *RoomController) List(ctx *gin.Context) {
	var rooms []entity.Room
	if err := c.roomService.List(&rooms); err != nil {
		newInternalErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (c *RoomController) Create(ctx *gin.Context) {
	var req createRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	room := entity.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	}

	if err := c.roomService.Create(&room); err != nil {
		newInternalErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"room": room})
}

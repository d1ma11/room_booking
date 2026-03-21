package controller

import (
	"net/http"
	"test-backend-1-d1ma11/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type SlotController struct {
	slotService service.SlotService
}

func NewSlotController(slotService service.SlotService) *SlotController {
	return &SlotController{slotService: slotService}
}

func (c *SlotController) List(ctx *gin.Context) {
	roomID := ctx.Param("roomId")
	dateStr := ctx.Query("date")

	if dateStr == "" {
		newErrorResponse(
			ctx,
			http.StatusBadRequest,
			service.NewError(service.ErrorType.InvalidRequest, "missing required query parameter: date"),
		)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		newErrorResponse(
			ctx,
			http.StatusBadRequest,
			service.NewError(service.ErrorType.InvalidRequest, "invalid date format, expected YYYY-MM-DD"),
		)
		return
	}

	slots, err := c.slotService.ListAvailableSlots(roomID, date)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"slots": slots})
}

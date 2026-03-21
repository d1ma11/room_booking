package controller

import (
	"fmt"
	log "log/slog"
	"net/http"
	"strconv"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingService service.BookingService
}

func NewBookingController(bookingService service.BookingService) *BookingController {
	return &BookingController{bookingService: bookingService}
}

type createBookingRequest struct {
	SlotId               string `json:"slotId" binding:"required"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

func (c *BookingController) BookSlot(ctx *gin.Context) {
	var req createBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErrorResponse(
			ctx,
			http.StatusBadRequest,
			service.NewError(service.ErrorType.InvalidRequest, err.Error()))
		return
	}
	log.Info(fmt.Sprintf("Request to book slot: %v", req))

	userId, _ := ctx.Get("user_id")
	booking := entity.Booking{
		UserID: userId.(string),
		SlotID: req.SlotId,
	}

	log.Info(fmt.Sprintf("Book room without conference link"))
	err := c.bookingService.Create(&booking, req.CreateConferenceLink)
	if err != nil {
		log.Info(fmt.Sprintf("There is an error while creating booking for slot. Error: %v", err))
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"booking": booking})
}

func (c *BookingController) CancelBooking(ctx *gin.Context) {
	bookingId := ctx.Param("bookingId")
	userId, _ := ctx.Get("user_id")

	var booking entity.Booking
	err := c.bookingService.CancelBooking(&booking, bookingId, userId.(string))
	if err != nil {
		log.Info(fmt.Sprintf("There is an error while canceling booking for slot. Error: %v", err))
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"booking": booking})
}

func (c *BookingController) ListMy(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")

	var bookings []entity.Booking
	if err := c.bookingService.ListByUserId(&bookings, userId.(string)); err != nil {
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func (c *BookingController) ListAdmin(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		newErrorResponse(ctx, http.StatusBadRequest, service.NewError(service.ErrorType.InvalidRequest, "invalid page parameter"))
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		newErrorResponse(ctx, http.StatusBadRequest, service.NewError(service.ErrorType.InvalidRequest, "pageSize must be between 1 and 100"))
		return
	}

	var bookings []entity.Booking
	var total int64
	err = c.bookingService.ListByAdminId(&bookings, page, pageSize, &total)
	if err != nil {
		log.Info(fmt.Sprintf("There is an error while listing booking for slot. Error: %v", err))
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}

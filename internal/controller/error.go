package controller

import (
	"errors"
	"net/http"
	"test-backend-1-d1ma11/internal/service"

	"github.com/gin-gonic/gin"
)

func newErrorResponse(ctx *gin.Context, errorStatus int, response error) {
	ctx.JSON(errorStatus, gin.H{"error": response})
}

func newInternalErrorResponse(ctx *gin.Context, response error) {
	ctx.JSON(http.StatusInternalServerError, response)
}

func handleServiceError(ctx *gin.Context, err error) {
	var appErr *service.ErrorResponse
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case service.ErrorType.SlotNotFound, service.ErrorType.BookingNotFound, service.ErrorType.RoomNotFound:
			newErrorResponse(ctx, http.StatusNotFound, appErr)
		case service.ErrorType.SlotAlreadyBooked, service.ErrorType.ScheduleExists:
			newErrorResponse(ctx, http.StatusConflict, appErr)
		case service.ErrorType.InvalidRequest:
			newErrorResponse(ctx, http.StatusBadRequest, appErr)
		case service.ErrorType.Forbidden:
			newErrorResponse(ctx, http.StatusForbidden, appErr)
		default:
			newInternalErrorResponse(ctx, err)
		}
		return
	}
	newInternalErrorResponse(ctx, err)
}

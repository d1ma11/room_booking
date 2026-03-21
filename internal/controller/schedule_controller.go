package controller

import (
	"fmt"
	log "log/slog"
	"net/http"
	"regexp"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type ScheduleController struct {
	scheduleService service.ScheduleService
}

func NewScheduleController(scheduleService service.ScheduleService) *ScheduleController {
	return &ScheduleController{scheduleService: scheduleService}
}

type createScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek" binding:"required,dive,min=1,max=7"`
	StartTime  string `json:"startTime" binding:"required"`
	EndTime    string `json:"endTime" binding:"required"`
}

var timePattern = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)

func (c *ScheduleController) Create(ctx *gin.Context) {
	roomId := ctx.Param("roomId")

	var req createScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		newErrorResponse(
			ctx,
			http.StatusBadRequest,
			service.NewError(service.ErrorType.InvalidRequest, err.Error()),
		)
		return
	}

	if !timePattern.MatchString(req.StartTime) || !timePattern.MatchString(req.EndTime) {
		newErrorResponse(
			ctx,
			http.StatusBadRequest,
			service.NewError(service.ErrorType.InvalidRequest, "time must be in HH:MM format"),
		)
		return
	}

	start, _ := time.Parse("15:04", req.StartTime)
	end, _ := time.Parse("15:04", req.EndTime)
	if !start.Before(end) {
		newErrorResponse(
			ctx,
			http.StatusBadRequest,
			service.NewError(service.ErrorType.InvalidRequest, "startTime must be before endTime"),
		)
		return
	}

	days := make(pq.Int32Array, len(req.DaysOfWeek))
	for i, d := range req.DaysOfWeek {
		days[i] = int32(d)
	}

	schedule := entity.Schedule{
		RoomID:     roomId,
		DaysOfWeek: days,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	err := c.scheduleService.Create(&schedule, roomId)
	if err != nil {
		log.Info(fmt.Sprintf("There is an error while creating schedule for room. Error: %v", err))
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"schedule": schedule})
}

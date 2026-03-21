package service

import (
	"fmt"
)

type ErrorTypes struct {
	InvalidRequest    string
	Unauthorized      string
	NotFound          string
	RoomNotFound      string
	SlotNotFound      string
	SlotAlreadyBooked string
	BookingNotFound   string
	Forbidden         string
	ScheduleExists    string
	InternalError     string
}

var ErrorType = ErrorTypes{
	InvalidRequest:    "INVALID_REQUEST",
	Unauthorized:      "UNAUTHORIZED",
	NotFound:          "NOT_FOUND",
	RoomNotFound:      "ROOM_NOT_FOUND",
	SlotNotFound:      "SLOT_NOT_FOUND",
	SlotAlreadyBooked: "SLOT_ALREADY_BOOKED",
	BookingNotFound:   "BOOKING_NOT_FOUND",
	Forbidden:         "FORBIDDEN",
	ScheduleExists:    "SCHEDULE_EXISTS",
	InternalError:     "INTERNAL_ERROR",
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Code %s: %s", e.Code, e.Message)
}

func NewError(errorCode, errorMessage string) *ErrorResponse {
	return &ErrorResponse{Code: errorCode, Message: errorMessage}
}

type InternalErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *InternalErrorResponse) Error() string {
	return fmt.Sprintf("Code %s: %s", e.Code, e.Message)
}

func newInternalError(errorMessage string) *InternalErrorResponse {
	return &InternalErrorResponse{Code: ErrorType.InternalError, Message: errorMessage}
}

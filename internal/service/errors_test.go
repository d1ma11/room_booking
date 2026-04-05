package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse_JSONTags(t *testing.T) {
	err := NewError(ErrorType.InvalidRequest, "test message")

	jsonBytes, marshalErr := json.Marshal(err)
	assert.NoError(t, marshalErr)

	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, `"code":"INVALID_REQUEST"`)
	assert.Contains(t, jsonStr, `"message":"test message"`)
	assert.NotContains(t, jsonStr, `"Code"`)
	assert.NotContains(t, jsonStr, `"Message"`)
}

func TestErrorResponse_ErrorMethod(t *testing.T) {
	err := NewError(ErrorType.RoomNotFound, "room not found")

	assert.Equal(t, "Code ROOM_NOT_FOUND: room not found", err.Error())
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("database timeout")

	assert.Equal(t, ErrorType.InternalError, err.Code)
	assert.Equal(t, "database timeout", err.Message)
	assert.Equal(t, "Code INTERNAL_ERROR: database timeout", err.Error())
}

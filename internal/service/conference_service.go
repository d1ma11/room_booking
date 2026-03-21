package service

import (
	"fmt"
	"time"
)

type ConferenceServiceImpl struct {
}

func NewConferenceService() *ConferenceServiceImpl {
	return &ConferenceServiceImpl{}
}

func (s *ConferenceServiceImpl) CreateConferenceLink(bookingID string) (string, error) {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		link, err := s.callExternalService(bookingID)
		if err == nil {
			return link, nil
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	return "", fmt.Errorf("failed to create conference link after %d retries: %w", maxRetries, lastErr)
}

func (s *ConferenceServiceImpl) callExternalService(bookingID string) (string, error) {
	return fmt.Sprintf("https://conference.com/meeting/%s", bookingID), nil
}

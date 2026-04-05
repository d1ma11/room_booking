package services

import (
	"test-backend-1-d1ma11/internal/service/auth"
	"time"

	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (m *JwtServiceMock) GenerateToken(userId, role, secret string, ttl time.Duration) (string, error) {
	args := m.Called(userId, role, secret, ttl)
	return args.String(0), args.Error(1)
}

func (m *JwtServiceMock) ParseToken(tokenStr, secret string) (*auth.Claims, error) {
	args := m.Called(tokenStr, secret)
	return args.Get(0).(*auth.Claims), args.Error(1)
}

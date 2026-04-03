package service

import (
	"errors"
	"test-backend-1-d1ma11/configs"
	"test-backend-1-d1ma11/internal/auth"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	FixedAdminUUID = "00000000-0000-0000-0000-000000000001"
	FixedUserUUID  = "00000000-0000-0000-0000-000000000002"
	emptyString    = ""
	userRole       = "user"
	adminRole      = "admin"
)

type UserServiceImpl struct {
	cfg      *configs.Config
	userRepo repository.UserRepository
}

func NewUserServiceImpl(cfg *configs.Config, userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{cfg: cfg, userRepo: userRepo}
}

func (s *UserServiceImpl) DummyLogin(role string) (string, error) {
	userId := FixedUserUUID
	if role == adminRole {
		userId = FixedAdminUUID
	}

	token, err := auth.GenerateToken(userId, role, s.cfg.JWT.Secret, s.cfg.JWT.TTL)
	if err != nil {
		return emptyString, newInternalError("failed to generate token")
	}

	return token, nil
}

func (s *UserServiceImpl) Register(email, password, role string) (entity.User, error) {
	var user entity.User
	if err := s.userRepo.GetByEmail(&user, email); err == nil {
		return entity.User{}, NewError(ErrorType.UserExists, "user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return entity.User{}, err
	}

	user = entity.User{
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
	if err = s.userRepo.Create(&user); err != nil {
		return entity.User{}, NewInternalError("failed to save user")
	}

	return user, nil
}

func (s *UserServiceImpl) Login(email, password string) (string, error) {
	var existingUser entity.User
	if err := s.userRepo.GetByEmail(&existingUser, email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return emptyString, NewError(ErrorType.InvalidCredentials, "invalid credentials")
		}
		return emptyString, NewInternalError("failed to find user")
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		return emptyString, NewError(ErrorType.InvalidCredentials, "invalid credentials")
	}

	token, err := auth.GenerateToken(existingUser.ID, existingUser.Role, s.cfg.JWT.Secret, s.cfg.JWT.TTL)
	if err != nil {
		return emptyString, newInternalError("failed to generate token")
	}

	return token, nil
}

package service

import (
	"errors"
	"test-backend-1-d1ma11/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	userService = NewUserServiceImpl(testConfig, mockUserRepo, mockJwtService)
}

func TestUserService_DummyLogin_Admin(t *testing.T) {
	// Arrange
	token := "testToken"

	mockJwtService.On("GenerateToken", fixedAdminUUID, adminRole, testConfig.JWT.Secret, testConfig.JWT.TTL).
		Return(token, nil).
		Once()

	// Act
	token, err := userService.DummyLogin(adminRole)

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_DummyLogin_User(t *testing.T) {
	// Arrange
	token := "testToken"

	mockJwtService.On("GenerateToken", fixedUserUUID, userRole, testConfig.JWT.Secret, testConfig.JWT.TTL).
		Return(token, nil).
		Once()

	// Act
	token, err := userService.DummyLogin(userRole)

	// Assert
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_DummyLogin_FailedGeneration(t *testing.T) {
	// Arrange
	mockJwtService.On("GenerateToken", fixedUserUUID, userRole, testConfig.JWT.Secret, testConfig.JWT.TTL).
		Return(emptyString, errors.New("generation failed")).
		Once()

	// Act
	token, err := userService.DummyLogin(userRole)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	var appErr *InternalErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.InternalError, appErr.Code)
	assert.Empty(t, token)
}

func TestUserService_Register_UserExists(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"

	mockUserRepo.On("GetByEmail", &entity.User{}, testEmail).
		Return(nil).
		Once()

	// Act
	user, err := userService.Register(testEmail, testPassword, userRole)

	// Assert
	assert.Empty(t, user)
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	var appErr *ErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.UserExists, appErr.Code)
	mockUserRepo.AssertNotCalled(t, "Create")
}

func TestUserService_Register_UserNotExists(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"

	mockUserRepo.On("GetByEmail", &entity.User{}, testEmail).
		Return(gorm.ErrRecordNotFound).
		Once()
	mockUserRepo.On("Create", mock.AnythingOfType("*entity.User")).
		Return(nil).
		Once()

	// Act
	user, err := userService.Register(testEmail, testPassword, userRole)

	// Assert
	assert.NotEmpty(t, user)
	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "Create", mock.AnythingOfType("*entity.User"))
}

func TestUserService_Register_FailedToSave(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"
	dbErr := errors.New("failed to save user")

	mockUserRepo.On("GetByEmail", &entity.User{}, testEmail).
		Return(gorm.ErrRecordNotFound).
		Once()
	mockUserRepo.On("Create", mock.AnythingOfType("*entity.User")).
		Return(dbErr).
		Once()

	// Act
	user, err := userService.Register(testEmail, testPassword, userRole)

	// Assert
	assert.Empty(t, user)
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	var appErr *InternalErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.InternalError, appErr.Code)
	mockUserRepo.AssertCalled(t, "Create", mock.AnythingOfType("*entity.User"))
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"

	mockUserRepo.On("GetByEmail", &entity.User{}, testEmail).
		Return(gorm.ErrRecordNotFound).
		Once()

	// Act
	token, err := userService.Login(testEmail, testPassword)

	assert.Empty(t, token)
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	var appErr *ErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.InvalidCredentials, appErr.Code)
	mockJwtService.AssertNotCalled(t, "GenerateToken")
}

func TestUserService_Login_DbProblem(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"
	dbErr := errors.New("internal db problem")

	mockUserRepo.On("GetByEmail", &entity.User{}, testEmail).
		Return(dbErr).
		Once()

	// Act
	token, err := userService.Login(testEmail, testPassword)

	assert.Empty(t, token)
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	var appErr *InternalErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.InternalError, appErr.Code)
	mockJwtService.AssertNotCalled(t, "GenerateToken")
}

func TestUserService_Login_FailedGeneration(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"
	jwtErr := errors.New("internal jwt problem")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	existingUser := &entity.User{
		ID:       fixedUserUUID,
		Email:    testEmail,
		Password: string(hashedPassword),
		Role:     userRole,
	}

	mockUserRepo.On("GetByEmail", mock.AnythingOfType("*entity.User"), testEmail).
		Return(nil).
		Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*entity.User)
			*userArg = *existingUser
		}).
		Once()
	mockJwtService.On("GenerateToken", fixedUserUUID, userRole, testConfig.JWT.Secret, testConfig.JWT.TTL).
		Return(emptyString, jwtErr).
		Once()

	// Act
	token, err := userService.Login(testEmail, testPassword)

	assert.Empty(t, token)
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	var appErr *InternalErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.InternalError, appErr.Code)
	assert.Equal(t, "failed to generate token", appErr.Message)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"
	badPassword := "badPassword"
	testToken := "testToken"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	existingUser := &entity.User{
		ID:       fixedUserUUID,
		Email:    testEmail,
		Password: string(hashedPassword),
		Role:     userRole,
	}

	mockUserRepo.On("GetByEmail", mock.AnythingOfType("*entity.User"), testEmail).
		Return(nil).
		Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*entity.User)
			*userArg = *existingUser
		}).
		Once()
	mockJwtService.On("GenerateToken", fixedUserUUID, userRole, testConfig.JWT.Secret, testConfig.JWT.TTL).
		Return(testToken, nil).
		Once()

	// Act
	token, err := userService.Login(testEmail, badPassword)

	assert.Empty(t, token)
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	var appErr *ErrorResponse
	errors.As(err, &appErr)
	assert.Equal(t, ErrorType.InvalidCredentials, appErr.Code)
	assert.Equal(t, "invalid credentials", appErr.Message)
}

func TestUserService_Login_Success(t *testing.T) {
	// Arrange
	testEmail := "testEmail"
	testPassword := "testPassword"
	testToken := "testToken"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	existingUser := &entity.User{
		ID:       fixedUserUUID,
		Email:    testEmail,
		Password: string(hashedPassword),
		Role:     userRole,
	}

	mockUserRepo.On("GetByEmail", mock.AnythingOfType("*entity.User"), testEmail).
		Return(nil).
		Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*entity.User)
			*userArg = *existingUser
		}).
		Once()
	mockJwtService.On("GenerateToken", fixedUserUUID, userRole, testConfig.JWT.Secret, testConfig.JWT.TTL).
		Return(testToken, nil).
		Once()

	// Act
	token, err := userService.Login(testEmail, testPassword)

	assert.NotEmpty(t, token)
	assert.Nil(t, err)
}

package controller

import (
	"fmt"
	log "log/slog"
	"net/http"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

type dummyLoginRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (c *UserController) DummyLogin(ctx *gin.Context) {
	var req dummyLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "invalid request",
		}})
		return
	}

	token, err := c.userService.DummyLogin(req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "failed to generate token",
		}})
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse{Token: token})
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type registerResponse struct {
	User entity.User `json:"user"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "invalid request",
		}})
		return
	}

	// pass email and password to user_service
	// getting token
	user, err := c.userService.Register(req.Email, req.Password, req.Role)
	if err != nil {
		newInternalErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, registerResponse{User: user})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "invalid request",
		}})
		return
	}

	token, err := c.userService.Login(req.Email, req.Password)
	if err != nil {
		log.Info(fmt.Sprintf("There is an error while canceling booking for slot. Error: %v", err))
		handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse{Token: token})
}

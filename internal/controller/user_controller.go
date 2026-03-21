package controller

import (
	"net/http"
	"test-backend-1-d1ma11/configs"
	"test-backend-1-d1ma11/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	FixedAdminUUID = "00000000-0000-0000-0000-000000000001"
	FixedUserUUID  = "00000000-0000-0000-0000-000000000002"
)

type AuthController struct {
	cfg *configs.Config
}

func NewAuthController(cfg *configs.Config) *AuthController {
	return &AuthController{cfg: cfg}
}

type dummyLoginRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (c *AuthController) DummyLogin(ctx *gin.Context) {
	var req dummyLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "invalid request",
		}})
		return
	}

	userId := FixedUserUUID
	if req.Role == "admin" {
		userId = FixedAdminUUID
	}

	token, err := auth.GenerateToken(userId, req.Role, c.cfg.JWT.Secret, c.cfg.JWT.TTL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "failed to generate token",
		}})
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse{Token: token})
}

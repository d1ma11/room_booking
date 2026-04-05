package middleware

import (
	"net/http"
	"strings"
	"test-backend-1-d1ma11/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.JwtService, secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "missing authorization header",
			}})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "invalid authorization header format",
			}})
			ctx.Abort()
			return
		}

		claims, err := jwtService.ParseToken(parts[1], secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "invalid or expired token",
			}})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != requiredRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "access denied",
			}})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

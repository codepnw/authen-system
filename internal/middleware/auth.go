package middleware

import (
	"net/http"
	"strings"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/utils/security"
	"github.com/gin-gonic/gin"
)

const UserContextKey = "user"

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	tokenCfg := security.NewJWTToken(cfg)

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "header is missing"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := tokenCfg.VerifyAccessToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		ctx.Set(UserContextKey, user)
		ctx.Next()
	}
}

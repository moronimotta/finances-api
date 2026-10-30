package server

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func sharedBearerToken(environmentVariable string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		expected := strings.TrimSpace(os.Getenv(environmentVariable))
		if expected == "" {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "API authorization is not configured"})
			ctx.Abort()
			return
		}
		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		provided := ""
		if strings.HasPrefix(authorization, "Bearer ") {
			provided = strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		}
		if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func adminToken() gin.HandlerFunc {
	return sharedBearerToken("PAYMENT_API_ADMIN_TOKEN")
}

func barberServiceToken() gin.HandlerFunc {
	return sharedBearerToken("BARBER_PAYMENT_API_TOKEN")
}

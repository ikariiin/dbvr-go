package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ikariiin/dbvr-go/utils"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := utils.ValidateToken(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

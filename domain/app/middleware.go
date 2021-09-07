package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func AuthMiddleware(findActiveProduct func(token string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			zap.S().Info("BearerAuthMiddleware error: empty token")
			c.AbortWithError(http.StatusUnauthorized, errors.New("empty token"))
			return
		}

		found := findActiveProduct(token)
		if !found {
			zap.S().Info("AuthMiddleware error: token was not found")
			c.AbortWithError(http.StatusUnauthorized, errors.New("bad token"))
			return
		}

		//c.Set(KeyUserID, payload.Subject)
		c.Next()
	}
}
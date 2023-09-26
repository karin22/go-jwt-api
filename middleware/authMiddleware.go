package middleware

import (
	"net/http"
	"strings"

	"helloGo/jwt-api/service"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

var hmacSampleSecret []byte

func AuthorizationMiddleware(c *gin.Context) {

	header := c.Request.Header.Get("Authorization")

	if header == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Provided token is invalid",
		})
		return
	}

	tokenString := strings.TrimPrefix(header, "Bearer ")
	if _, err := service.ValidateToken(tokenString); err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
		return
	}

}

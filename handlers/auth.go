package handlers

import (
	"fmt"
	"net/http"

	"atlan-proto/config"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || token != fmt.Sprintf("Bearer %s", config.SecretKey) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}
	c.Next()
}

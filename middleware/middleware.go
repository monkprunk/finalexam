package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginMiddleware(c *gin.Context) {
	log.Println("Start Middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()

}

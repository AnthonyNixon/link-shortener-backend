package auth_handler

import (
	"github.com/anthonynixon/link-shortener-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var TEMP_PASSWORD = os.Getenv("TEMP_PASS")

func AddAuthV1(router *gin.Engine) {
	router.POST("/token", CreateJWT)
}

func CreateJWT(c *gin.Context) {
	if c.GetHeader("password") != TEMP_PASSWORD {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	token, err := auth.New("anthony@nixon.dev")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

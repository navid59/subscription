package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAuthorized(c *gin.Context) {

	
	token := c.Request.Header.Get("token")
	if len(token) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Countinu with JWT",
			"token":   token,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Token didn't find in header!",
		})
	}
}

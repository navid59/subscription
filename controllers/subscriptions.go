package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Subscription struct {
	UserName     string    `json:"userName" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	LastName     string    `json:"lastName" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Adress       string    `json:"adress"`
	Tel          string    `json:"tel"`
	MrchanUserID string    `json:"mrchanUserID" binding:"required"`
	Plan         string    `json:"plan"  binding:"required"`
	StartDate    time.Time `json:"startDate" binding:"required"`
	EndDate      time.Time `json:"endDate"`
}

func SetSubscription(c *gin.Context) {
	var newSubscription Subscription
	if err := c.BindJSON(&newSubscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	if err := valid(newSubscription.Email); err != true {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Email validation failed!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "subscription ADD - POST",
		"details": "Ok Temporary",
	})
}

func Unsubscribe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscribe a Member - POST",
	})
}

func GetSubscriptionList(c *gin.Context) {
	// ListSubscription()
	c.JSON(http.StatusOK, gin.H{
		"message": "get subscription List - POST",
	})
}

func GetSubscription(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get subscription info - POST",
	})
}

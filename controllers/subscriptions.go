package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"fmt"
	"context"
	"log"
	"cloud.google.com/go/datastore"
	"crypto/sha256"
	"encoding/base64"
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

type Subscriber struct {
	PlanId       string `json:"planId"`
	UserId       string `json:"UserId" binding:"required"`
	MrchanUserID string `json:"mrchanUserID"`
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

	/** 
		add the new subscriber 
	*/
	ctx := context.Background()
	projectID := "netopia-payments"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
				   
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	newSubscriber := &newSubscription
	strHash := sha256.Sum256([]byte(newSubscriber.Password))
	newSubscriber.Password = base64.StdEncoding.EncodeToString(strHash[:])
	
	key := datastore.IncompleteKey("subscribers", nil)
	key.Namespace = "recurring"
	if _, err := client.Put(ctx, key, newSubscriber); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":    "ERROR",
			"message": "Error durring ADD Data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "subscriber registered successfully!",
	})
}

func Unsubscribe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscribe a Member - POST",
	})
}

func GetSubscriptionList(c *gin.Context) {
	ListSubscription(c)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "get subscription List - POST",
	// })
}

func GetSubscriptionStatus(c *gin.Context) {
	var member Subscriber
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get subscription info - POST",
	})
}

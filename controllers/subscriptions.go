package controllers

import (
	"cloud.google.com/go/datastore"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"log"
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
	id           *datastore.Key
}

type Subscriber struct {
	PlanId       string `json:"planId"`
	UserEmail    string `json:"userEmail"`
	MrchanUserID string `json:"mrchanUserID" binding:"required"`
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
	/* Validate Request Body & asssign to VAR member*/
	var member Subscriber
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), member.MrchanUserID); !ok {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	// Creates a client.
	ctx := context.Background()
	projectID := "netopia-payments"
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	/* Get Example of search */
	query := datastore.NewQuery("subscribers").Namespace("recurring").
		Filter("Plan =", member.PlanId).
		Filter("Email =", member.UserEmail).
		Filter("MrchanUserID =", member.MrchanUserID)

	var entitis []Subscription
	it := client.Run(ctx, query)
	for {
		var entity Subscription
		_, err := it.Next(&entity)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next task: %v", err)
		} else {
			entitis = append(entitis, entity)
			fmt.Printf("Nume : %q, Prenume:  %q\n", entity.Name, entity.LastName)
		}
	}

	/* Get Example */
	// subscriberkey := datastore.IDKey("subscribers", 5752521540763648, nil)
	// subscriberkey.Namespace = "recurring"
	// entity := &Subscription{}
	// if err := client.Get(ctx, subscriberkey, entity); err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"code":    http.StatusNotFound,
	// 		"message": "Data not found",
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"members": entitis,
	})
}

func GetLastSubscriberInfo() {
	/* return following info */
	// "lastTransactionAt" ,  "lastTransactionStatus"
}

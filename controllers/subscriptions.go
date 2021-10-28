package controllers

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"time"
	// "crypto/sha256"
	// "encoding/base64"
)


type Subscription struct {
	id           *datastore.Key
	Name         string    `json:"Name" binding:"required"`
	LastName     string    `json:"LastName" binding:"required"`
	Email        string    `json:"Email" binding:"required"`
	Adress       string    `json:"Adress"`
	Tel          string    `json:"Tel"`
	MrchanUserID string    `json:"MrchanUserID" binding:"required"`
	Plan         string    `json:"Plan"  binding:"required"`
	StartDate    time.Time `json:"StartDate" binding:"required"`
	EndDate      time.Time `json:"EndDate"`
	Status       bool    
	Flags        string    
	CreatedAt    time.Time 
	UpdatedAt    string
}

type Subscriber struct {
	PlanId       string `json:"planId" binding:"required"`
	UserEmail    string `json:"userEmail" binding:"required"`
	MrchanUserID string `json:"mrchanUserID" binding:"required"`
}

type User struct {
	UserId       int64  `json:"userId" binding:"required"`
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

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), newSubscription.MrchanUserID); !ok {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
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

	/* To make hash string like password, if there is case*/
	// strHash := sha256.Sum256([]byte(newSubscriber.Password))
	// newSubscriber.Password = base64.StdEncoding.EncodeToString(strHash[:])

	// set predefined data
	newSubscriber.Status = true
	newSubscriber.Flags = "subscribed"
	newSubscriber.CreatedAt = time.Now()
	// newSubscriber.UpdatedAt = ""

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

func GetSubscriptionList(c *gin.Context) {
	ListSubscription(c)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "get subscription List - POST",
	// })
}

func GetSubscriptionSearch(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"members": entitis,
	})
}

func GetSubscriptionStatus(c *gin.Context) {
	/* Validate Request Body & asssign to VAR member*/
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), user.MrchanUserID); !ok {
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

	/* Get User Info */
	subscriberkey := datastore.IDKey("subscribers", user.UserId, nil)
	subscriberkey.Namespace = "recurring"
	entity := &Subscription{}
	if err := client.Get(ctx, subscriberkey, entity); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Data not found",
		})
		return
	}

	// Verify Audience & Ownership of data
	if ownership := entity.MrchanUserID; len(ownership) > 0 && ownership != user.MrchanUserID {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"member": entity,
	})
}

func Unsubscribe(c *gin.Context) {
	/* Validate Request Body & asssign to VAR member*/
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
			"error":   err,
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), user.MrchanUserID); !ok {
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

	/* Get User Info */
	subscriberkey := datastore.IDKey("subscribers", user.UserId, nil)
	subscriberkey.Namespace = "recurring"
	entity := &Subscription{}
	if err := client.Get(ctx, subscriberkey, entity); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Data not found",
		})
		return
	}

	// Verify Audience & Ownership of data
	if ownership := entity.MrchanUserID; len(ownership) > 0 && ownership != user.MrchanUserID {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	if entity.Status == false {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusNotModified,
			"message": "Member already is unsubscribed",
		})
		return
	}  

	/* Unscubscribe the Member */
	entity.Status = false
	entity.Flags = "unsubscribed"
	tmpTime := time.Now()
	entity.UpdatedAt = tmpTime.String()

	if _, err := client.Put(ctx, subscriberkey, entity); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Request failed!",
		})
	return
    }

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Unsubscribed successfully!",
	})
}

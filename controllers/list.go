/**
Make a list of all users
Note : Currently is all the data, without Filter

Todo :
	- Not able to costumize json (return less than all existan properties in kind),
		ex. password shouldn't return
	- Most be filter by MrchanUserID
*/
package controllers

import (
	"context"
	// "fmt"
	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type StrcMember struct {
	UserName     string    `datastore:"userName"`
	Password     string    `datastore:"password"`
	Name         string    `datastore:"name"`
	LastName     string    `datastore:"lastName"`
	Email        string    `datastore:"email"`
	Adress       string    `datastore:"adress"`
	Tel          string    `datastore:"tel"`
	MrchanUserID string    `datastore:"mrchanUserID"`
	Plan         string    `datastore:"plan"`
	StartDate    time.Time `datastore:"startDate"`
	EndDate      time.Time `datastore:"endDate"`
	PaymentInfo  string    `datastore:"paymentInfo"`
	Status       string    `datastore:"status"`
	Flags        string    `datastore:"flags"`
	CreatedAt    time.Time `datastore:"createdAt"`
	UpdatedAt    time.Time `datastore:"updatedAt"`
	Id           int64     // The integer ID used in the datastore.
}

type Merchant struct {
	MerchantID string `json:"merchantID" binding:"required"`
}

func ListSubscription(c *gin.Context) {

	var merchant Merchant
	if err := c.BindJSON(&merchant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request Errors",
		})
		return
	}

	// Verify Audience
	if _, ok := IsLicensed(c.Request.Header.Get("token"), merchant.MerchantID); !ok {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Request has conflict with your authion",
		})
		return
	}

	// Creates a client.
	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	projectID := "netopia-payments"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var members []StrcMember
	query := datastore.NewQuery("subscribers").Namespace("recurring").
	Filter("MrchanUserID =", merchant.MerchantID)

	keys, err := client.GetAll(ctx, query, &members); 
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": http.StatusUnprocessableEntity, "message": "", "error": err})
		return
	}
	for i, key := range keys {
		members[i].Id = key.ID
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
	})

}

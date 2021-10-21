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
	// id           int64     // The integer ID used in the datastore.
}

func ListSubscription(c *gin.Context) {
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
	query := datastore.NewQuery("subscribers").Namespace("recurring")
	if _, err := client.GetAll(ctx, query, &members); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": http.StatusUnprocessableEntity, "message": "", "error": err})
		return
	}

	// Temporary - to see the result in a loop
	// for i, key := range keys {
	// 	Subscribers[i].id = key.ID
	// 	fmt.Println(Subscribers[i].UserName)
	// }

	c.JSON(http.StatusOK, gin.H{
		"members": members,
	})

}

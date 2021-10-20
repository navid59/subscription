package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

type SubscriptionDB struct {
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
	id           int64     // The integer ID used in the datastore.
}

func ListSubscription() {
	ctx := context.Background()

	// Set your Google Cloud Platform project ID.
	projectID := "netopia-payments"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var Subscribers []*SubscriptionDB
	query := datastore.NewQuery("subscribers").Namespace("recurring")
	keys, err := client.GetAll(ctx, query, &Subscribers)
	if err != nil {
		panic(err)
	}

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		Subscribers[i].id = key.ID
		fmt.Println(Subscribers[i].UserName)
	}
}
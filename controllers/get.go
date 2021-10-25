package controllers

import (
	// "context"
	"fmt"
	// "cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	// "log"
	"net/http"
	// "time"
)


func GetMethod(c *gin.Context) {
	fmt.Println("\n'GetMethod' called")
	IdValue := c.Params.ByName("IdValue")
	message := "GetMethod Called With Param: " + IdValue
	c.JSON(http.StatusOK, message)
  
	ReqPayload := make([]byte, 1024)
	ReqPayload, err := c.GetRawData()
	if err != nil {
		  fmt.Println(err)
		  return
	}
	fmt.Println("Request Payload Data: ", string(ReqPayload))
  }

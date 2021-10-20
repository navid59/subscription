package main

import (
	"github.com/gin-gonic/gin"
	"recurring/controllers"
)

func main() {
	r := gin.Default()
	r.GET("/", controllers.Home)
	r.NoRoute(controllers.PageNotFound)
	r.POST("/subscription", controllers.SetSubscription)
	r.POST("/subscription/cancel", controllers.Unsubscribe)
	r.POST("/subscription/list", controllers.GetSubscriptionList)
	r.POST("/subscription/status", controllers.GetSubscription)
	r.HandleMethodNotAllowed = true
	r.Run(":8080")
}

package main

import (
	"github.com/gin-gonic/gin"
	"recurring/controllers"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger()) // Logger middleware 

	// Create Sub Router for  customised API version
	authorized := r.Group("/api/v1")
	authorized.Use(controllers.IsAuthorized)
	authorized.GET("/:IdValue", controllers.GetMethod)
	authorized.GET("/jwt", controllers.Welcome)

	r.GET("/", controllers.Home)
	r.NoRoute(controllers.PageNotFound)
	r.POST("/subscription", controllers.SetSubscription)
	r.POST("/subscription/cancel", controllers.Unsubscribe)
	r.POST("/subscription/list", controllers.GetSubscriptionList)
	r.POST("/subscription/status", controllers.GetSubscriptionStatus)
	r.HandleMethodNotAllowed = true
	r.Run(":8080")
}

func middleware() {

}
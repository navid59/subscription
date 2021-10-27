package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("mySuperSecrtePhrase")

func IsAuthorized(c *gin.Context) {

	token := c.Request.Header.Get("token")
	if len(token) > 0 {

		_, err := jwt.Parse(c.Request.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": err.Error(),
			})
		}

		/* Just keep as ex - temporary */
		// c.JSON(http.StatusOK, gin.H{
		// 	"code":    http.StatusOK,
		// 	"message": "Compare Audience - TEMPORARY",
		// 	"token":   token, // Was define at parser part
		// 	"V1":      c.Query("order"),
		// 	"param":   c.Params.ByName("IdValue"),
		// })

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Missing token!",
		})
	}
}

func GetTokenAud(c *gin.Context) {

}

func IsLicensed(tokenStr, MrchanUserID string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if _, ok := claims["aud"]; ok && claims["aud"] == MrchanUserID {
			return claims, true
		}

		return claims, false
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = "securing@token096"
var salt = "yat@12#!*"
var pepper = "pep_887!@"

type Claims struct {
	UserId primitive.ObjectID `json:"userId"`
	jwt.RegisteredClaims
}

func (c *controller) signToken(ctx *gin.Context, userId primitive.ObjectID) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := base64.StdEncoding.EncodeToString([]byte(jwtKey))
	jwtToken = base64.StdEncoding.EncodeToString([]byte(pepper + jwtToken))
	jwtToken = base64.StdEncoding.EncodeToString([]byte(jwtToken + salt))

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(jwtToken))
	if err != nil{
		fmt.Println("Error generating token: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"msg": "Something went wrong at Server Side. Please try again later. If issue still persist, kindly contact our administrator.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg": tokenString,
	})
}

func (c *controller) validateToken(ctx *gin.Context, token string) primitive.ObjectID {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		jwtToken := base64.StdEncoding.EncodeToString([]byte(jwtKey))
		jwtToken = base64.StdEncoding.EncodeToString([]byte(pepper + jwtToken))
		jwtToken = base64.StdEncoding.EncodeToString([]byte(jwtToken + salt))
		
		return []byte(jwtToken), nil
	})

	if err != nil {
		fmt.Printf("Error while validating token: %v\n", err.Error())
		if err == jwt.ErrSignatureInvalid{
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": 1005,
				"msg": "Unauthorized Access Denied",
			})
			return primitive.NilObjectID
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 1006,
			"msg": "Invalid Request",
		})
		return primitive.NilObjectID
	}

	if !tkn.Valid{
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": 1005,
			"msg": "Unauthorized Access Denied",
		})
		return primitive.NilObjectID
	}

	return claims.UserId
}


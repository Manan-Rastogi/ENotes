package controllers

import (
	"fmt"
	"net/http"

	"github.com/Manan-Rastogi/enotes/configs"
	"github.com/Manan-Rastogi/enotes/models"
	"github.com/Manan-Rastogi/enotes/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller interface {
	SignUpController(ctx *gin.Context)
	LogInController(ctx *gin.Context)
	SaveNotesController(ctx *gin.Context)

	signToken(ctx *gin.Context, userId primitive.ObjectID)
	validateToken(ctx *gin.Context, token string) primitive.ObjectID
}

type controller struct {
	Service services.Service
}

func NewController(s services.Service) Controller {
	return &controller{
		Service: s,
	}
}

func (c *controller) SignUpController(ctx *gin.Context) {
	var user models.User

	Mongodb, _ := ctx.Get("mongodb")
	mongodb := Mongodb.(models.Mongo)

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": 1001,
			"msg":    "Invalid Input",
			"error":  err.Error(),
		})
		return
	}

	status, msg, err := c.Service.SignUpService(mongodb, user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": status,
			"msg":    msg,
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": status,
		"msg":    msg,
	})
}

func (c *controller) LogInController(ctx *gin.Context) {

	Mongodb, _ := ctx.Get("mongodb")
	mongodb := Mongodb.(models.Mongo)

	inputs, err := ctx.MultipartForm()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": 1001,
			"msg":    "Invalid Input",
			"error":  err.Error(),
		})
		return
	}

	email := inputs.Value["email"][0]
	password := inputs.Value["password"][0]

	id, err := c.Service.LogInService(mongodb, email, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": 1004,
			"msg":    err.Error(),
		})
		return
	}

	c.signToken(ctx, id)

}

func (c *controller) SaveNotesController(ctx *gin.Context) {
	var notes models.Notes
	Mongodb, _ := ctx.Get("mongodb")
	mongodb := Mongodb.(models.Mongo)

	if err := ctx.ShouldBindJSON(&notes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": 1001,
			"msg":    "Invalid Input",
			"error":  err.Error(),
		})
		return
	}

	token := ctx.GetHeader("authorization")[len(configs.BEARER_SCHEMA):]
	fmt.Printf("token: %v\n", token)
	userId := c.validateToken(ctx, token)

	if !userId.IsZero() {
		code, status, msg := c.Service.SaveNotesService(mongodb, userId, notes)

		ctx.JSON(code, gin.H{
			"status": status,
			"msg":    msg,
		})
	}

}

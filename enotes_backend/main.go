package main

import (
	"net/http"

	"github.com/Manan-Rastogi/enotes/controllers"
	"github.com/Manan-Rastogi/enotes/models"
	"github.com/Manan-Rastogi/enotes/services"
	"github.com/gin-gonic/gin"
)

var (
	service services.Service = services.NewService()
	controller controllers.Controller = controllers.NewController(service)
	mongodb = models.Connect()
)

func main() {
	
	// defer func() {
	// 	if err := mongodb.Client.Disconnect(mongodb.Ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	router := gin.New()

	
	router.Use()

	api := router.Group("api/v1")
	api.GET("/welcome", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "Welcome To E-Notes",
		})
	})

	api.POST("signup", SetContext(), controller.SignUpController)
	api.POST("login", SetContext(), controller.LogInController)
	api.POST("save_notes", SetContext(), controller.SaveNotesController)
	

	router.Run(":8081")

}

func SetContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("mongodb", mongodb)
	}
}
package api

import (
	"go_avito/internal/handlers"
	"go_avito/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	storage.InitDB()

	r := gin.Default()
	r.POST("/team/add", handlers.AddTeam)
	r.GET("/team/get", handlers.GetTeam)
	r.POST("/users/setIsActive", handlers.SetIsActive)
	r.POST("/pullRequest/create", handlers.NewPR)
	r.POST("/pullRequest/merge", handlers.MergePR)
	r.POST("/pullRequest/reassign", handlers.Reassign)
	r.GET("/users/getReview", handlers.GetReview)

	r.Run(":8080")

	log.Println("Server down")
}

package main

import (
	"log"

	"github.com/BochkaDeyalo/jsonplaceholder.api/config"
	"github.com/BochkaDeyalo/jsonplaceholder.api/controller"
	"github.com/BochkaDeyalo/jsonplaceholder.api/logger"
	"github.com/BochkaDeyalo/jsonplaceholder.api/service"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.Load(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logger.NewSlogLogger(config.LogLevel)
	logger.Info("Starting application", "logLevel", config.LogLevel)
	service := service.NewService(config)
	controller := controller.NewController(service, logger)
	r := gin.Default()

	r.GET("/posts", controller.GetPosts)
	r.GET("/posts/:id", controller.GetPostByID)
	r.POST("/posts", controller.CreatePost)
	r.PUT("/posts/:id", controller.UpdatePost)
	r.PATCH("/posts/:id", controller.PatchPost)
	r.DELETE("/posts/:id", controller.DeletePost)

	if err := r.Run(); err != nil {
		logger.Error("failed to run server: %v", err)
	}
}

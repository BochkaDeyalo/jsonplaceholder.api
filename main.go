package main

import (
	"log"

	"github.com/BochkaDeyalo/jsonplaceholder.api/config"
	"github.com/BochkaDeyalo/jsonplaceholder.api/controller"
	_ "github.com/BochkaDeyalo/jsonplaceholder.api/docs"
	"github.com/BochkaDeyalo/jsonplaceholder.api/logger"
	"github.com/BochkaDeyalo/jsonplaceholder.api/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           JSONPlaceholder API
// @version         1.0
// @description     A REST API for managing posts using JSONPlaceholder external service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes http https
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(); err != nil {
		logger.Error("failed to run server: %v", err)
	}
}

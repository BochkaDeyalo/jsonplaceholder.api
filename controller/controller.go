package controller

import (
	"net/http"

	apierror "github.com/BochkaDeyalo/jsonplaceholder.api/error"
	"github.com/BochkaDeyalo/jsonplaceholder.api/logger"
	"github.com/BochkaDeyalo/jsonplaceholder.api/model"
	"github.com/BochkaDeyalo/jsonplaceholder.api/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Logger    logger.Logger
	ApiCaller service.ApiCaller
	validator *validator.Validate
}

func NewController(apiCaller service.ApiCaller, log logger.Logger) *Controller {
	return &Controller{
		ApiCaller: apiCaller,
		validator: validator.New(),
		Logger:    log,
	}
}

func (pc *Controller) CreatePost(c *gin.Context) {
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		pc.Logger.Warn("Invalid JSON", "error", err.Error())
		apierror.SendError(c, apierror.BadRequest, err.Error())
		return
	}

	pc.Logger.Debug("Creating post", "userId", request.UserID, "title", request.Title)

	if err := pc.validator.Struct(request); err != nil {
		pc.Logger.Warn("Validation failed", "error", err.Error())
		apierror.SendError(c, apierror.ValidationFailed, err.Error())
		return
	}

	result, err := pc.ApiCaller.CreatePost(request)
	if err != nil {
		pc.Logger.Error("External service error", "error", err.Error())
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	pc.Logger.Info("Post created successfully", "postId", result.ID)
	c.JSON(http.StatusCreated, result)

}

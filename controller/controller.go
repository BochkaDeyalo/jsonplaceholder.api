package controller

import (
	"net/http"

	apierror "github.com/BochkaDeyalo/jsonplaceholder.api/error"
	"github.com/BochkaDeyalo/jsonplaceholder.api/model"
	"github.com/BochkaDeyalo/jsonplaceholder.api/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	ApiCaller service.ApiCaller
	validator *validator.Validate
}

func NewController(apiCaller service.ApiCaller) *Controller {
	return &Controller{
		ApiCaller: apiCaller,
		validator: validator.New(),
	}
}

func (pc *Controller) CreatePost(c *gin.Context) {
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		apierror.SendError(c, apierror.BadRequest, err.Error())
		return
	}

	if err := pc.validator.Struct(request); err != nil {
		apierror.SendError(c, apierror.ValidationFailed, err.Error())
		return
	}

	result, err := pc.ApiCaller.CreatePost(request)
	if err != nil {
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

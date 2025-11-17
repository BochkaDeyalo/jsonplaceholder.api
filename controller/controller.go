package controller

import (
	"net/http"

	"github.com/BochkaDeyalo/jsonplaceholder.api/model"
	"github.com/BochkaDeyalo/jsonplaceholder.api/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *service.Service
}

func NewController(Service *service.Service) *Controller {
	return &Controller{
		Service: Service,
	}
}

func (pc *Controller) CreatePost(c *gin.Context) {
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := pc.Service.CreatePost(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

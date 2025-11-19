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

func (pc *Controller) GetPosts(c *gin.Context) {
	result, err := pc.ApiCaller.GetPosts()
	if err != nil {
		pc.Logger.Error("External service error", "error", err.Error())
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	pc.Logger.Info("GET posts successfully", "postsCount", len(result.Posts))
	c.JSON(http.StatusOK, result)
}

func (pc *Controller) GetPostByID(c *gin.Context) {
	id := c.Param("id")

	pc.Logger.Debug("Getting post by ID", "id", id)

	result, err := pc.ApiCaller.GetPostByID(id)
	if err != nil {
		pc.Logger.Error("External service error", "error", err.Error())
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	pc.Logger.Info("GET post by ID successfully", "postId", result.ID)
	c.JSON(http.StatusOK, result)
}

func (pc *Controller) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		pc.Logger.Warn("Invalid JSON", "error", err.Error())
		apierror.SendError(c, apierror.BadRequest, err.Error())
		return
	}

	pc.Logger.Debug("Updating post", "id", id, "userId", request.UserID, "title", request.Title)

	if err := pc.validator.Struct(request); err != nil {
		pc.Logger.Warn("Validation failed", "error", err.Error())
		apierror.SendError(c, apierror.ValidationFailed, err.Error())
		return
	}

	result, err := pc.ApiCaller.UpdatePost(id, request)
	if err != nil {
		pc.Logger.Error("External service error", "error", err.Error())
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	pc.Logger.Info("Post updated successfully", "postId", result.ID)
	c.JSON(http.StatusOK, result)
}

func (pc *Controller) PatchPost(c *gin.Context) {
	id := c.Param("id")

	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		pc.Logger.Warn("Invalid JSON", "error", err.Error())
		apierror.SendError(c, apierror.BadRequest, err.Error())
		return
	}

	pc.Logger.Debug("Patching post", "id", id)

	result, err := pc.ApiCaller.PatchPost(id, request)
	if err != nil {
		pc.Logger.Error("External service error", "error", err.Error())
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	pc.Logger.Info("Post patched successfully", "postId", result.ID)
	c.JSON(http.StatusOK, result)
}

func (pc *Controller) DeletePost(c *gin.Context) {
	id := c.Param("id")

	pc.Logger.Debug("Deleting post", "id", id)

	err := pc.ApiCaller.DeletePost(id)
	if err != nil {
		pc.Logger.Error("External service error", "error", err.Error())
		apierror.SendError(c, apierror.ExternalServiceError, err.Error())
		return
	}

	pc.Logger.Info("Post deleted successfully", "id", id)
	c.Data(http.StatusNoContent, "application/json", nil)
}

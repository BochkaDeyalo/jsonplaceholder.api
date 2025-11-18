package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BochkaDeyalo/jsonplaceholder.api/controller"
	apierror "github.com/BochkaDeyalo/jsonplaceholder.api/error"
	"github.com/BochkaDeyalo/jsonplaceholder.api/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockApiCaller struct {
	mock.Mock
}

func (m *MockApiCaller) CreatePost(request model.RequestBody) (*model.ResponseBody, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResponseBody), args.Error(1)
}

func TestCreatePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	expectedResponse := &model.ResponseBody{
		UserID: 1,
		ID:     101,
		Title:  "Test Post",
		Body:   "Test Body",
	}

	mockApiCaller.On("CreatePost", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.UserID == 1 && req.Title == "Test Post" && req.Body == "Test Body"
	})).Return(expectedResponse, nil)

	ctrl := controller.NewController(mockApiCaller)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "Test Post",
		Body:   "Test Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreatePost(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.ResponseBody
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Title, response.Title)
	assert.Equal(t, expectedResponse.Body, response.Body)

	mockApiCaller.AssertExpectations(t)
}

func TestCreatePost_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	ctrl := controller.NewController(mockApiCaller)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "",
		Body:   "Test Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreatePost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Validation failed", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "CreatePost")
}

func TestCreatePost_ValidationFailed_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	ctrl := controller.NewController(mockApiCaller)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "Test Title",
		Body:   "",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreatePost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Validation failed", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "CreatePost")
}

func TestCreatePost_ValidationFailed_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	ctrl := controller.NewController(mockApiCaller)

	requestBody := model.RequestBody{
		UserID: 0,
		Title:  "Test Title",
		Body:   "Test Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreatePost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Validation failed", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "CreatePost")
}

func TestCreatePost_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	serviceError := errors.New("external service unavailable")
	mockApiCaller.On("CreatePost", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.UserID == 1 && req.Title == "Test Post" && req.Body == "Test Body"
	})).Return(nil, serviceError)

	ctrl := controller.NewController(mockApiCaller)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "Test Post",
		Body:   "Test Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreatePost(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Equal(t, "External service error", response.Message)
	assert.Equal(t, serviceError.Error(), response.Details)

	mockApiCaller.AssertExpectations(t)
}

func TestCreatePost_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	ctrl := controller.NewController(mockApiCaller)

	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreatePost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Invalid request", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "CreatePost")
}

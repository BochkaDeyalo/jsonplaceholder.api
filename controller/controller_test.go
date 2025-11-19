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

func (m *MockApiCaller) GetPosts() (*model.ResponseArrayBody, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResponseArrayBody), args.Error(1)
}

func (m *MockApiCaller) GetPostByID(id string) (*model.ResponseBody, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResponseBody), args.Error(1)
}

func (m *MockApiCaller) UpdatePost(id string, request model.RequestBody) (*model.ResponseBody, error) {
	args := m.Called(id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResponseBody), args.Error(1)
}

func (m *MockApiCaller) PatchPost(id string, request model.RequestBody) (*model.ResponseBody, error) {
	args := m.Called(id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResponseBody), args.Error(1)
}

func (m *MockApiCaller) DeletePost(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...any) {
	m.Called(msg, args)
}

func TestCreatePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	expectedResponse := &model.ResponseBody{
		UserID: 1,
		ID:     101,
		Title:  "Test Post",
		Body:   "Test Body",
	}

	mockApiCaller.On("CreatePost", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.UserID == 1 && req.Title == "Test Post" && req.Body == "Test Body"
	})).Return(expectedResponse, nil)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

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
	mockLogger := new(MockLogger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

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
	mockLogger := new(MockLogger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

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
	mockLogger := new(MockLogger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

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
	mockLogger := new(MockLogger)
	serviceError := errors.New("external service unavailable")
	mockApiCaller.On("CreatePost", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.UserID == 1 && req.Title == "Test Post" && req.Body == "Test Body"
	})).Return(nil, serviceError)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

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
	mockLogger := new(MockLogger)

	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

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

func TestGetPosts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	expectedResponse := &model.ResponseArrayBody{
		Posts: []model.ResponseBody{
			{
				UserID: 1,
				ID:     1,
				Title:  "Test Post 1",
				Body:   "Test Body 1",
			},
			{
				UserID: 2,
				ID:     2,
				Title:  "Test Post 2",
				Body:   "Test Body 2",
			},
		},
	}

	mockApiCaller.On("GetPosts").Return(expectedResponse, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetPosts(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.ResponseArrayBody
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedResponse.Posts), len(response.Posts))
	assert.Equal(t, expectedResponse.Posts[0].UserID, response.Posts[0].UserID)
	assert.Equal(t, expectedResponse.Posts[0].ID, response.Posts[0].ID)
	assert.Equal(t, expectedResponse.Posts[0].Title, response.Posts[0].Title)
	assert.Equal(t, expectedResponse.Posts[0].Body, response.Posts[0].Body)
	assert.Equal(t, expectedResponse.Posts[1].UserID, response.Posts[1].UserID)
	assert.Equal(t, expectedResponse.Posts[1].ID, response.Posts[1].ID)
	assert.Equal(t, expectedResponse.Posts[1].Title, response.Posts[1].Title)
	assert.Equal(t, expectedResponse.Posts[1].Body, response.Posts[1].Body)

	mockApiCaller.AssertExpectations(t)
}

func TestGetPosts_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	serviceError := errors.New("external service unavailable")

	mockApiCaller.On("GetPosts").Return(nil, serviceError)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetPosts(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Equal(t, "External service error", response.Message)
	assert.Equal(t, serviceError.Error(), response.Details)

	mockApiCaller.AssertExpectations(t)
}

func TestGetPostByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	expectedResponse := &model.ResponseBody{
		UserID: 1,
		ID:     1,
		Title:  "Test Post",
		Body:   "Test Body",
	}

	mockApiCaller.On("GetPostByID", "1").Return(expectedResponse, nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/post/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.GetPostByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.ResponseBody
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Title, response.Title)
	assert.Equal(t, expectedResponse.Body, response.Body)

	mockApiCaller.AssertExpectations(t)
}

func TestGetPostByID_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	serviceError := errors.New("external service unavailable")

	mockApiCaller.On("GetPostByID", "1").Return(nil, serviceError)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodGet, "/post/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.GetPostByID(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Equal(t, "External service error", response.Message)
	assert.Equal(t, serviceError.Error(), response.Details)

	mockApiCaller.AssertExpectations(t)
}

func TestUpdatePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	expectedResponse := &model.ResponseBody{
		UserID: 1,
		ID:     1,
		Title:  "Updated Post",
		Body:   "Updated Body",
	}

	mockApiCaller.On("UpdatePost", "1", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.UserID == 1 && req.Title == "Updated Post" && req.Body == "Updated Body"
	})).Return(expectedResponse, nil)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "Updated Post",
		Body:   "Updated Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/post/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.UpdatePost(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.ResponseBody
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Title, response.Title)
	assert.Equal(t, expectedResponse.Body, response.Body)

	mockApiCaller.AssertExpectations(t)
}

func TestUpdatePost_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "",
		Body:   "Updated Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/post/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.UpdatePost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Validation failed", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "UpdatePost")
}

func TestUpdatePost_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)

	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodPut, "/post/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.UpdatePost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Invalid request", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "UpdatePost")
}

func TestUpdatePost_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	serviceError := errors.New("external service unavailable")

	mockApiCaller.On("UpdatePost", "1", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.UserID == 1 && req.Title == "Updated Post" && req.Body == "Updated Body"
	})).Return(nil, serviceError)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	requestBody := model.RequestBody{
		UserID: 1,
		Title:  "Updated Post",
		Body:   "Updated Body",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/post/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.UpdatePost(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Equal(t, "External service error", response.Message)
	assert.Equal(t, serviceError.Error(), response.Details)

	mockApiCaller.AssertExpectations(t)
}

func TestPatchPost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	expectedResponse := &model.ResponseBody{
		UserID: 1,
		ID:     1,
		Title:  "Patched Title",
		Body:   "Original Body",
	}

	mockApiCaller.On("PatchPost", "1", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.Title == "Patched Title"
	})).Return(expectedResponse, nil)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	requestBody := model.RequestBody{
		Title: "Patched Title",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/posts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.PatchPost(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.ResponseBody
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.UserID, response.UserID)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Title, response.Title)
	assert.Equal(t, expectedResponse.Body, response.Body)

	mockApiCaller.AssertExpectations(t)
}

func TestPatchPost_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)

	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodPatch, "/posts/1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.PatchPost(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Invalid request", response.Message)
	assert.NotEmpty(t, response.Details)

	mockApiCaller.AssertNotCalled(t, "PatchPost")
}

func TestPatchPost_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	serviceError := errors.New("external service unavailable")

	mockApiCaller.On("PatchPost", "1", mock.MatchedBy(func(req model.RequestBody) bool {
		return req.Title == "Patched Title"
	})).Return(nil, serviceError)

	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	requestBody := model.RequestBody{
		Title: "Patched Title",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/posts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.PatchPost(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Equal(t, "External service error", response.Message)
	assert.Equal(t, serviceError.Error(), response.Details)

	mockApiCaller.AssertExpectations(t)
}

func TestDeletePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)

	mockApiCaller.On("DeletePost", "1").Return(nil)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.DeletePost(c)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockApiCaller.AssertExpectations(t)
}

func TestDeletePost_ExternalServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockApiCaller := new(MockApiCaller)
	mockLogger := new(MockLogger)
	serviceError := errors.New("external service unavailable")

	mockApiCaller.On("DeletePost", "1").Return(serviceError)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	ctrl := controller.NewController(mockApiCaller, mockLogger)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	ctrl.DeletePost(c)

	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response apierror.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, response.Code)
	assert.Equal(t, "External service error", response.Message)
	assert.Equal(t, serviceError.Error(), response.Details)

	mockApiCaller.AssertExpectations(t)
}

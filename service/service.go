package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BochkaDeyalo/jsonplaceholder.api/config"
	"github.com/BochkaDeyalo/jsonplaceholder.api/model"
)

type ApiCaller interface {
	CreatePost(request model.RequestBody) (*model.ResponseBody, error)
	GetPosts() (*model.ResponseArrayBody, error)
	GetPostByID(id string) (*model.ResponseBody, error)
	UpdatePost(id string, request model.RequestBody) (*model.ResponseBody, error)
	PatchPost(id string, request model.RequestBody) (*model.ResponseBody, error)
	DeletePost(id string) error
}

type Service struct {
	client *http.Client
	config *config.Config
}

func NewService(cfg *config.Config) ApiCaller {
	return &Service{
		client: &http.Client{},
		config: cfg,
	}
}

func (s *Service) CreatePost(request model.RequestBody) (*model.ResponseBody, error) {
	apiURL := s.config.APIURL
	endpoint := "/posts"
	url := apiURL + endpoint

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response model.ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GetPosts() (*model.ResponseArrayBody, error) {
	apiURL := s.config.APIURL
	endpoint := "/posts"
	url := apiURL + endpoint

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response model.ResponseArrayBody
	if err := json.NewDecoder(resp.Body).Decode(&response.Posts); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GetPostByID(id string) (*model.ResponseBody, error) {
	apiURL := s.config.APIURL
	endpoint := "/posts/" + id
	url := apiURL + endpoint

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response model.ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) UpdatePost(id string, request model.RequestBody) (*model.ResponseBody, error) {
	apiURL := s.config.APIURL
	endpoint := "/posts/" + id
	url := apiURL + endpoint

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response model.ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) PatchPost(id string, request model.RequestBody) (*model.ResponseBody, error) {
	apiURL := s.config.APIURL
	endpoint := "/posts/" + id
	url := apiURL + endpoint

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response model.ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) DeletePost(id string) error {
	apiURL := s.config.APIURL
	endpoint := "/posts/" + id
	url := apiURL + endpoint

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/BochkaDeyalo/jsonplaceholder.api/config"
	"github.com/BochkaDeyalo/jsonplaceholder.api/model"
)

type Service struct {
	client *http.Client
	config *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		client: &http.Client{},
		config: cfg,
	}
}

func (s *Service) CreatePost(request model.RequestBody) (*model.RequestBody, error) {
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

	return &request, nil
}

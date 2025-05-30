package services

import (
	"encoding/json"
	"net/http"
	"testing"
	"todo-app/tests/helpers"

	"github.com/stretchr/testify/assert"
)

// UserService - Service Layer untuk User operations
type UserService struct {
	httpClient *helpers.TestHTTPClient
	t          *testing.T
}

// UserResponse represents user response
type UserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// UserResult contains user operation result
type UserResult struct {
	Response   *UserResponse
	StatusCode int
	Success    bool
}

// NewUserService creates new user service
func NewUserService(httpClient *helpers.TestHTTPClient, t *testing.T) *UserService {
	return &UserService{
		httpClient: httpClient,
		t:          t,
	}
}

// GetProfile gets user profile
func (service *UserService) GetProfile(cookies []*http.Cookie) *UserResult {
	httpResp := service.httpClient.GetWithCookies("/user/profile", cookies)

	var userResp UserResponse
	err := json.Unmarshal(httpResp.Body, &userResp)
	assert.NoError(service.t, err, "Should be able to parse user profile response")

	return &UserResult{
		Response:   &userResp,
		StatusCode: httpResp.StatusCode,
		Success:    httpResp.StatusCode == http.StatusOK && userResp.Success,
	}
}

// GetProfileWithoutAuth tries to get profile without authentication
func (service *UserService) GetProfileWithoutAuth() *UserResult {
	httpResp := service.httpClient.Get("/user/profile")

	// For unauthorized access, we might not get valid JSON
	var userResp UserResponse
	json.Unmarshal(httpResp.Body, &userResp) // Don't assert error here as it might fail for 401

	return &UserResult{
		Response:   &userResp,
		StatusCode: httpResp.StatusCode,
		Success:    false, // Should always be false for unauthorized access
	}
}

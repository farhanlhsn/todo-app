package services

import (
	"encoding/json"
	"net/http"
	"testing"
	"todo-app/tests/helpers"

	"github.com/stretchr/testify/assert"
)

// AuthService - Service Layer untuk Authentication
type AuthService struct {
	httpClient *helpers.TestHTTPClient
	t          *testing.T
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents register request
type RegisterRequest struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResult contains authentication result
type AuthResult struct {
	Response *AuthResponse
	Cookies  []*http.Cookie
	Success  bool
}

// NewAuthService creates new auth service
func NewAuthService(httpClient *helpers.TestHTTPClient, t *testing.T) *AuthService {
	return &AuthService{
		httpClient: httpClient,
		t:          t,
	}
}

// Register performs user registration
func (service *AuthService) Register(nama, email, password string) *AuthResult {
	request := RegisterRequest{
		Nama:     nama,
		Email:    email,
		Password: password,
	}

	httpResp := service.httpClient.Post("/auth/register", request)

	var authResp AuthResponse
	err := json.Unmarshal(httpResp.Body, &authResp)
	assert.NoError(service.t, err, "Should be able to parse register response")

	return &AuthResult{
		Response: &authResp,
		Cookies:  httpResp.Cookies,
		Success:  httpResp.StatusCode == http.StatusOK && authResp.Success,
	}
}

// Login performs user login
func (service *AuthService) Login(email, password string) *AuthResult {
	request := LoginRequest{
		Email:    email,
		Password: password,
	}

	httpResp := service.httpClient.Post("/auth/login", request)

	var authResp AuthResponse
	err := json.Unmarshal(httpResp.Body, &authResp)
	assert.NoError(service.t, err, "Should be able to parse login response")

	return &AuthResult{
		Response: &authResp,
		Cookies:  httpResp.Cookies,
		Success:  httpResp.StatusCode == http.StatusOK && authResp.Success,
	}
}

// Logout performs user logout
func (service *AuthService) Logout(cookies []*http.Cookie) *AuthResult {
	httpResp := service.httpClient.PostWithCookies("/auth/logout", nil, cookies)

	var authResp AuthResponse
	err := json.Unmarshal(httpResp.Body, &authResp)
	assert.NoError(service.t, err, "Should be able to parse logout response")

	return &AuthResult{
		Response: &authResp,
		Cookies:  httpResp.Cookies,
		Success:  httpResp.StatusCode == http.StatusOK && authResp.Success,
	}
}

// GetAuthCookie extracts Authorization cookie from cookies
func (service *AuthService) GetAuthCookie(cookies []*http.Cookie) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == "Authorization" {
			return cookie
		}
	}
	return nil
}

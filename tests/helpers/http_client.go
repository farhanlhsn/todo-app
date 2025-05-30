package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHTTPClient - HTTP Client Layer
type TestHTTPClient struct {
	router *gin.Engine
	t      *testing.T
}

// HTTPResponse represents standard HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Cookies    []*http.Cookie
}

// NewTestHTTPClient creates new HTTP client for testing
func NewTestHTTPClient(router *gin.Engine, t *testing.T) *TestHTTPClient {
	return &TestHTTPClient{
		router: router,
		t:      t,
	}
}

// Post makes POST request
func (client *TestHTTPClient) Post(path string, payload interface{}) *HTTPResponse {
	return client.makeRequest("POST", path, payload, nil)
}

// PostWithCookies makes POST request with cookies
func (client *TestHTTPClient) PostWithCookies(path string, payload interface{}, cookies []*http.Cookie) *HTTPResponse {
	return client.makeRequest("POST", path, payload, cookies)
}

// Get makes GET request
func (client *TestHTTPClient) Get(path string) *HTTPResponse {
	return client.makeRequest("GET", path, nil, nil)
}

// GetWithCookies makes GET request with cookies
func (client *TestHTTPClient) GetWithCookies(path string, cookies []*http.Cookie) *HTTPResponse {
	return client.makeRequest("GET", path, nil, cookies)
}

// makeRequest - internal method for making HTTP requests
func (client *TestHTTPClient) makeRequest(method, path string, payload interface{}, cookies []*http.Cookie) *HTTPResponse {
	var reqBody *bytes.Buffer

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		assert.NoError(client.t, err, "Should be able to marshal request payload")
		reqBody = bytes.NewBuffer(jsonData)
	}

	var req *http.Request
	if reqBody != nil {
		req = httptest.NewRequest(method, path, reqBody)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add cookies if provided
	if cookies != nil {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	w := httptest.NewRecorder()
	client.router.ServeHTTP(w, req)

	return &HTTPResponse{
		StatusCode: w.Code,
		Body:       w.Body.Bytes(),
		Cookies:    w.Result().Cookies(),
	}
}

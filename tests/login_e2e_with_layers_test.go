package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo-app/controllers"
	"todo-app/initializers"
	"todo-app/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Simple response structure for testing
type TestResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Login request structure for testing
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register request structure for testing
type RegisterRequest struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TestLoginE2E_WithLayeringPattern - E2E Test dengan Layering Pattern
func TestLoginE2E_WithLayeringPattern(t *testing.T) {
	// Setup Gin dalam mode testing
	gin.SetMode(gin.TestMode)

	// Load environment variables dari .env file
	initializers.LoadEnvVariables()

	// Initialize aplikasi dengan database dari .env
	setupTestEnvironment()

	// Setup router dengan routes yang akan ditest
	router := setupTestRouter()

	// Bersihkan database sebelum test
	cleanupDatabase()

	t.Run("Layer_1_Test_Registration_With_Database", func(t *testing.T) {
		// Test registration dengan data valid
		validData := RegisterRequest{
			Nama:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(validData)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return success
		assert.Equal(t, http.StatusOK, w.Code, "Valid registration should succeed")

		var response TestResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should be able to parse response")
		assert.True(t, response.Success, "Registration should be successful")
		assert.Equal(t, "User registered successfully", response.Message)

		t.Logf("Layer 1 Complete: Registration with database successful")
	})

	t.Run("Layer_2_Test_Login_With_Database", func(t *testing.T) {
		// Test login dengan kredensial yang valid
		loginData := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return success
		assert.Equal(t, http.StatusOK, w.Code, "Valid login should succeed")

		var response TestResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should be able to parse response")
		assert.True(t, response.Success, "Login should be successful")
		assert.Equal(t, "Login successful", response.Message)

		// Verify JWT cookie is set
		cookies := w.Result().Cookies()
		var authCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "Authorization" {
				authCookie = cookie
				break
			}
		}
		assert.NotNil(t, authCookie, "Authorization cookie should be set")
		assert.NotEmpty(t, authCookie.Value, "JWT token should not be empty")

		t.Logf("Layer 2 Complete: Login with database successful")
	})

	t.Run("Layer_3_Test_Protected_Route_With_Auth", func(t *testing.T) {
		// Step 1: Register user dulu (supaya test independen)
		registerData := RegisterRequest{
			Nama:     "Test User Layer 3",
			Email:    "test3@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(registerData)
		regReq := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
		regReq.Header.Set("Content-Type", "application/json")
		regW := httptest.NewRecorder()
		router.ServeHTTP(regW, regReq)

		// Pastikan registration berhasil
		assert.Equal(t, http.StatusOK, regW.Code, "Registration should succeed for Layer 3")

		// Step 2: Login untuk dapat token
		loginData := LoginRequest{
			Email:    "test3@example.com",
			Password: "password123",
		}

		loginJsonData, _ := json.Marshal(loginData)
		loginReq := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginJsonData))
		loginReq.Header.Set("Content-Type", "application/json")
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)

		// Debug: Cek response login
		assert.Equal(t, http.StatusOK, loginW.Code, "Login should succeed first")

		// Ambil cookie dari login response
		cookies := loginW.Result().Cookies()

		// Debug: Print semua cookies yang ada
		t.Logf("Debug: Found %d cookies", len(cookies))
		for i, cookie := range cookies {
			t.Logf("Debug: Cookie %d - Name: %s, Value: %s", i, cookie.Name, cookie.Value[:min(20, len(cookie.Value))]+"...")
		}

		var authCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "Authorization" {
				authCookie = cookie
				break
			}
		}

		// Null check untuk authCookie
		if authCookie == nil {
			t.Fatalf("Authorization cookie not found! Available cookies: %v", func() []string {
				var names []string
				for _, c := range cookies {
					names = append(names, c.Name)
				}
				return names
			}())
		}

		// Step 3: Test akses protected route dengan token
		profileReq := httptest.NewRequest("GET", "/user/profile", nil)
		profileReq.AddCookie(authCookie)
		profileW := httptest.NewRecorder()
		router.ServeHTTP(profileW, profileReq)

		// Should be able to access
		assert.Equal(t, http.StatusOK, profileW.Code, "Should be able to access protected route with valid token")

		var profileResponse TestResponse
		err := json.Unmarshal(profileW.Body.Bytes(), &profileResponse)
		assert.NoError(t, err, "Should be able to parse profile response")
		assert.True(t, profileResponse.Success, "Profile request should be successful")

		t.Logf("✅ Layer 3 Complete: Protected route access successful")
	})

	t.Run("Layer_4_Test_Middleware_Authentication", func(t *testing.T) {
		// Test protected route without auth
		req := httptest.NewRequest("GET", "/user/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return unauthorized
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Protected route should require authentication")

		t.Logf("Layer 4 Complete: Middleware authentication working correctly")
	})

	t.Run("Layer_5_Test_Invalid_Credentials", func(t *testing.T) {
		// Setup: Register user untuk test wrong password
		setupUser := RegisterRequest{
			Nama:     "Test User Layer 5",
			Email:    "test5@example.com",
			Password: "correctpassword",
		}

		jsonSetup, _ := json.Marshal(setupUser)
		setupReq := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonSetup))
		setupReq.Header.Set("Content-Type", "application/json")
		setupW := httptest.NewRecorder()
		router.ServeHTTP(setupW, setupReq)
		assert.Equal(t, http.StatusOK, setupW.Code, "Setup user should be registered")

		// Test 1: Email yang tidak ada sama sekali
		testCases := []struct {
			name            string
			email           string
			password        string
			expectedMessage string
		}{
			{
				name:            "Non-existent Email",
				email:           "notexist@example.com",
				password:        "password123",
				expectedMessage: "Invalid email or password",
			},
			{
				name:            "Wrong Password",
				email:           "test5@example.com",
				password:        "wrongpassword",
				expectedMessage: "Invalid email or password",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				loginData := LoginRequest{
					Email:    tc.email,
					Password: tc.password,
				}

				jsonData, _ := json.Marshal(loginData)
				req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should return unauthorized
				assert.Equal(t, http.StatusUnauthorized, w.Code, "Invalid credentials should be rejected")

				var response TestResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Should be able to parse response")
				assert.False(t, response.Success, "Invalid login should fail")
				assert.Contains(t, response.Message, tc.expectedMessage, "Should contain proper error message")
			})
		}

		t.Logf("Layer 5 Complete: Invalid credentials properly rejected")
	})

	// Cleanup setelah semua test selesai
	t.Cleanup(func() {
		cleanupDatabase()
		t.Logf("Database cleanup completed")
	})

	t.Logf("All E2E layering pattern tests with database completed successfully!")
}

// Helper function untuk setup test environment
func setupTestEnvironment() {
	testDbUrl := os.Getenv("TEST_DATABASE_URL")
	if testDbUrl != "" {
		os.Setenv("DATABASE_URL", testDbUrl)
	}

	// Connect to test database
	initializers.ConnectToDB()

	// Sync database tables
	initializers.SyncDatabase()
}

// Helper function untuk setup router
func setupTestRouter() *gin.Engine {
	router := gin.New()

	// Setup routes yang akan ditest
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.LogIn)
		auth.POST("/logout", middlewares.RequiredAuth, controllers.LogOut)
	}

	user := router.Group("/user")
	{
		user.Use(middlewares.RequiredAuth)
		user.GET("/profile", controllers.Profile)
	}

	return router
}

// Helper function untuk cleanup database
func cleanupDatabase() {
	if initializers.DB != nil {
		initializers.DB.Exec("DELETE FROM users")
		initializers.DB.Exec("DELETE FROM tasks")
		initializers.DB.Exec("DELETE FROM task_categories")
	}
}

// Helper function min untuk debug
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

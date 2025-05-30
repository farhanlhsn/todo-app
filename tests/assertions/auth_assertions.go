package assertions

import (
	"net/http"
	"testing"
	"todo-app/initializers"
	"todo-app/models"
	"todo-app/tests/services"

	"github.com/stretchr/testify/assert"
)

// AuthAssertions - Assertion Layer untuk Authentication
type AuthAssertions struct {
	t *testing.T
}

// NewAuthAssertions creates new auth assertions
func NewAuthAssertions(t *testing.T) *AuthAssertions {
	return &AuthAssertions{t: t}
}

// AssertLoginSuccess verifies successful login
func (a *AuthAssertions) AssertLoginSuccess(result *services.AuthResult) {
	assert.True(a.t, result.Success, "Login should be successful")
	assert.Equal(a.t, "Login successful", result.Response.Message)

	// Verify JWT cookie is set
	authCookie := a.GetAuthCookie(result.Cookies)
	assert.NotNil(a.t, authCookie, "Authorization cookie should be set")
	assert.NotEmpty(a.t, authCookie.Value, "JWT token should not be empty")
}

// AssertLoginFailure verifies failed login
func (a *AuthAssertions) AssertLoginFailure(result *services.AuthResult, expectedMessage string) {
	assert.False(a.t, result.Success, "Login should fail")
	assert.Contains(a.t, result.Response.Message, expectedMessage, "Should return expected error message")
}

// AssertRegistrationSuccess verifies successful registration
func (a *AuthAssertions) AssertRegistrationSuccess(result *services.AuthResult) {
	assert.True(a.t, result.Success, "Registration should be successful")
	assert.Equal(a.t, "User registered successfully", result.Response.Message)
}

// AssertLogoutSuccess verifies successful logout
func (a *AuthAssertions) AssertLogoutSuccess(result *services.AuthResult) {
	assert.True(a.t, result.Success, "Logout should be successful")
	assert.Equal(a.t, "User logged out successfully", result.Response.Message)
}

// AssertUserLoggedInDatabase verifies user login status in database
func (a *AuthAssertions) AssertUserLoggedInDatabase(email string, expectedStatus bool) {
	var user models.User
	err := initializers.DB.First(&user, "email = ?", email).Error
	assert.NoError(a.t, err, "Should be able to find user in database")
	assert.Equal(a.t, expectedStatus, user.IsLoggedIn, "User login status should match expected")
}

// AssertProfileAccessSuccess verifies successful profile access
func (a *AuthAssertions) AssertProfileAccessSuccess(result *services.UserResult) {
	assert.True(a.t, result.Success, "Profile access should be successful")
	assert.Equal(a.t, http.StatusOK, result.StatusCode, "Should return 200 OK")
}

// AssertProfileAccessDenied verifies denied profile access
func (a *AuthAssertions) AssertProfileAccessDenied(result *services.UserResult) {
	assert.False(a.t, result.Success, "Profile access should be denied")
	assert.Equal(a.t, http.StatusUnauthorized, result.StatusCode, "Should return 401 Unauthorized")
}

// GetAuthCookie helper to extract Authorization cookie
func (a *AuthAssertions) GetAuthCookie(cookies []*http.Cookie) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == "Authorization" {
			return cookie
		}
	}
	return nil
}

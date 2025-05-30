package controllers

import (
	"net/http"
	"os"
	"time"
	"todo-app/helpers"
	"todo-app/initializers"
	"todo-app/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Nama     string `json:"nama" binding:"required,min=2,max=100"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=100"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Invalid request: "+err.Error()))
		return
	}

	// Check if user already exists by email
	var existingUser models.User
	initializers.DB.First(&existingUser, "email = ?", body.Email)
	if existingUser.ID != 0 {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "User already exists with this email"))
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error hashing password"))
		return
	}

	// Create a new user
	user := models.User{
		Nama:     body.Nama,
		Email:    body.Email,
		Password: string(hash),
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error creating user", result.Error.Error()))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponseWithoutData(true, "User registered successfully"))
}

func LogIn(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Invalid request: "+err.Error()))
		return
	}

	// Check if user exists by email
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "Invalid email or password"))
		return
	}
	// If user is already logged in, return an error
	if user.IsLoggedIn {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User is already logged in"))
		return
	}

	// Compare the provided password with the hashed password in the database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "Invalid email or password"))
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	})

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "JWT Secret Key is missing"))
		return
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error generating token", err.Error()))
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	// Update user login status
	user.IsLoggedIn = true
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error updating user login status", err.Error()))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponseWithoutData(true, "Login successful"))
}

func LogOut(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Update user login status
	userModel.IsLoggedIn = false
	if err := initializers.DB.Save(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error updating user logout status", err.Error()))
		return
	}

	// Clear the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponseWithoutData(true, "User logged out successfully"))
}

func Profile(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Return user profile information
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "User profile retrieved successfully", helpers.FormatUserResponse(userModel)))
}

func Stats(c *gin.Context) {
	// Get the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	var taskCount int64
	initializers.DB.Model(&models.Task{}).Where("user_id = ?", userModel.ID).Count(&taskCount)

	var completedCount int64
	initializers.DB.Model(&models.Task{}).Where("user_id = ? AND is_completed = ?", userModel.ID, true).Count(&completedCount)

	var overdueCount int64
	initializers.DB.Model(&models.Task{}).Where("user_id = ? AND is_completed = ? AND due_date < ?", userModel.ID, false, time.Now()).Count(&overdueCount)

	stats := map[string]interface{}{
		"task_count":      taskCount,
		"completed_count": completedCount,
		"overdue_count":   overdueCount,
	}

	c.JSON(http.StatusOK, helpers.FormatResponse(true, (userModel.Nama+"'s statistics retrieved successfully"), stats))
}

package controllers

import (
	"fmt"
	"net/http"
	"time"
	"todo-app/helpers"
	"todo-app/initializers"
	"todo-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTask(c *gin.Context) {
	var body struct {
		Title       string `json:"Title" binding:"required,min=1,max=255"`
		Description string `json:"Description" binding:"max=1000"`
		DueDate     string `json:"DueDate,omitempty"`
		Category    string `json:"Category" binding:"max=100"`
		Priority    string `json:"Priority" binding:"oneof=none low medium high"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Invalid request: "+err.Error()))
		return
	}
	// Take user from context
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert userValue to models.User
	currentUser, ok := userValue.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}
	// Validate required fields
	if body.Title == "" {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Title is required"))
		return
	}
	if body.Description == "" {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Description is required"))
		return
	}
	//parse due date
	layoutFormat := "2006-01-02 15:04"
	date, _ := time.Parse(layoutFormat, body.DueDate)

	//search for category
	var category models.TaskCategory
	var categoryID uint = 0
	if body.Category != "" {
		result := initializers.DB.Where("name = ?", body.Category).First(&category)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error retrieving category", result.Error.Error()))
			return
		}
		if result.Statement != nil {
			// If category is found, set the CategoryID
			categoryID = category.ID
		}
		if result.RowsAffected == 0 {
			// Get available categories (global categories and user-specific categories)
			var availableCategories []string
			initializers.DB.Model(&models.TaskCategory{}).
				Where("user_id IS NULL OR user_id = ?", currentUser.ID).
				Pluck("name", &availableCategories)

			c.JSON(http.StatusBadRequest, gin.H{
				"success":              false,
				"message":              "Category not found",
				"available_categories": availableCategories,
			})
			return
		}
	}

	// Create a new task
	task := models.Task{
		Title:       body.Title,
		Description: body.Description,
		UserID:      currentUser.ID,
		Priority:    models.Priority(body.Priority),
		DueDate:     &date,
		CategoryID:  &categoryID,
		IsCompleted: false,
	}

	result := initializers.DB.Create(&task)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error creating task", result.Error.Error()))
		return
	}

	// Preload the category for the response
	initializers.DB.Preload("Category").First(&task, task.ID)

	c.JSON(http.StatusCreated, helpers.FormatResponse(true, "Task created successfully", helpers.FormatTaskResponse(task)))
}

func GetTasks(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find all tasks belonging to the current user
	var tasks []models.Task
	initializers.DB.Preload("Category").Where("user_id = ?", userModel.ID).Find(&tasks)

	// Return tasks list
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Tasks retrieved successfully", helpers.FormatTasksResponse(tasks)))
}

func GetTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find task by ID and ensure it belongs to current user
	var task models.Task
	result := initializers.DB.Preload("Category").Where("user_id = ?", userModel.ID).Where("id = ?", taskID).First(&task)

	// Check if task was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Task not found"))
		return
	}

	// Return task data
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Task retrieved successfully", helpers.FormatTaskResponse(task)))
}

func UpdateTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID := c.Param("id")

	// Parse request body
	var body struct {
		Title       string `json:"Title" binding:"required,min=1,max=255"`
		Description string `json:"Description" binding:"max=1000"`
		DueDate     string `json:"DueDate"`
		Category    string `json:"Category" binding:"max=100"`
		Priority    string `json:"Priority" binding:"oneof=none low medium high"`
		IsCompleted bool   `json:"IsCompleted"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Invalid request: "+err.Error()))
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find task by ID and ensure it belongs to current user
	var task models.Task
	result := initializers.DB.Preload("Category").Where("user_id = ?", userModel.ID).Where("id = ?", taskID).First(&task)

	// Check if task was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Task not found"))
		return
	}

	// Update task fields
	task.Title = body.Title
	task.Description = body.Description
	task.Priority = models.Priority(body.Priority)
	task.IsCompleted = body.IsCompleted
	// Parse due date if provided
	if body.DueDate != "" {
		layoutFormat := "2006-01-02 15:04"
		date, err := time.Parse(layoutFormat, body.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Invalid due date format"))
			return
		}
		task.DueDate = &date
	} else {
		// If no due date provided, keep the existing one

	}
	// Search for category if provided
	var category models.TaskCategory
	var categoryID *uint = nil
	if body.Category != "" {
		result := initializers.DB.Where("name = ?", body.Category).First(&category)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error retrieving category", result.Error.Error()))
			return
		}
		if result.RowsAffected > 0 {
			// If category is found, set the CategoryID
			categoryID = &category.ID
		} else {
			// If category not found, return available categories
			var availableCategories []string
			initializers.DB.Model(&models.TaskCategory{}).
				Where("user_id IS NULL OR user_id = ?", userModel.ID).
				Pluck("name", &availableCategories)
			c.JSON(http.StatusBadRequest, gin.H{
				"success":              false,
				"message":              "Category not found",
				"available_categories": availableCategories,
			})
			return
		}
	} else {
		// If no category provided, keep the existing one
		if task.CategoryID != nil {
			categoryID = task.CategoryID
		}
	}
	// Update the CategoryID field
	task.CategoryID = categoryID

	// Save changes to database
	initializers.DB.Save(&task)

	// Reload the task with category relationship
	initializers.DB.Preload("Category").First(&task, task.ID)

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Task updated successfully", helpers.FormatTaskResponse(task)))
}

func CompletedTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find task by ID and ensure it belongs to current user
	var task models.Task
	result := initializers.DB.Preload("Category").Where("user_id = ?", userModel.ID).Where("id = ?", taskID).First(&task)

	// Check if task was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Task not found"))
		return
	}

	// Mark task as completed by setting status to true
	task.IsCompleted = true
	initializers.DB.Save(&task)

	// Reload the task with category relationship
	initializers.DB.Preload("Category").First(&task, task.ID)

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Task marked as completed", helpers.FormatTaskResponse(task)))
}
func DeleteTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find task by ID and ensure it belongs to current user
	var task models.Task
	result := initializers.DB.Preload("Category").Where("user_id = ?", userModel.ID).Where("id = ?", taskID).First(&task)

	// Check if task was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Task not found"))
		return
	}

	// Perform soft delete (updates DeletedAt field)
	initializers.DB.Delete(&task)

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponseWithoutData(true, "Task deleted successfully"))
}

func GetCompletedTasks(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find all completed tasks belonging to the current user
	var tasks []models.Task
	initializers.DB.Preload("Category").Where("user_id = ? AND is_completed = ?", userModel.ID, true).Find(&tasks)

	// Return completed tasks list
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Completed tasks retrieved successfully", helpers.FormatTasksResponse(tasks)))
}

func UncompletedTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find task by ID (including soft-deleted tasks) and ensure it belongs to current user
	var task models.Task
	result := initializers.DB.Unscoped().Preload("Category").Where("user_id = ?", userModel.ID).Where("id = ?", taskID).First(&task)

	// Check if task was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Task not found"))
		return
	}

	// Mark task as uncompleted and restore if it was deleted
	task.IsCompleted = false
	task.DeletedAt = gorm.DeletedAt{} // Clear the DeletedAt field to restore the task
	initializers.DB.Save(&task)

	// Reload the task with category relationship
	initializers.DB.Preload("Category").First(&task, task.ID)

	// Return success response
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Task marked as uncompleted", helpers.FormatTaskResponse(task)))
}

func CreateCategory(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, helpers.FormatResponseWithoutData(false, "Invalid request"))
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Create a new task category
	category := models.TaskCategory{
		Name:   body.Name,
		UserID: &userModel.ID,
	}

	result := initializers.DB.Create(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error creating category", result.Error.Error()))
		return
	}

	// Preload the user for the response
	initializers.DB.Preload("User").First(&category, category.ID)
	c.JSON(http.StatusCreated, helpers.FormatResponse(true, "Category created successfully", helpers.FormatCategoryResponse(category)))
}

func GetCategories(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find all categories belonging to the current user
	var categories []models.TaskCategory
	initializers.DB.Where("user_id = ? OR user_id IS NULL", userModel.ID).Find(&categories)

	// Return categories list
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Categories retrieved successfully", helpers.FormatCategoriesResponse(categories)))
}

func GetTasksByCategory(c *gin.Context) {
	// Get category name from URL parameter
	categoryName := c.Param("name")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find category by name and ensure it belongs to current user
	var category models.TaskCategory
	result := initializers.DB.Where("name LIKE ?", "%"+categoryName+"%").
		Where("user_id = ? OR user_id IS NULL", userModel.ID).
		Order("user_id DESC"). // User categories first, then default
		First(&category)

	// Check if category was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Category not found"))
		return
	}

	// Find all tasks belonging to the current user and the specified category
	var tasks []models.Task
	initializers.DB.Preload("Category").Where("user_id = ? AND category_id = ?", userModel.ID, category.ID).Find(&tasks)

	// Return tasks list
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Tasks retrieved successfully", helpers.FormatTasksResponse(tasks)))
}

func DeleteCategory(c *gin.Context) {
	// Get category ID from URL parameter
	categoryID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find category by ID
	var category models.TaskCategory
	result := initializers.DB.Where("id = ?", categoryID).First(&category)

	// Check if category was found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, helpers.FormatResponseWithoutData(false, "Category not found"))
		return
	}

	// Check if this is a default category (user_id IS NULL)
	if category.UserID == nil {
		// This is a default/global category - only remove it from user's tasks
		// Update all user's tasks that use this category to have null category_id
		updateResult := initializers.DB.Model(&models.Task{}).
			Where("user_id = ? AND category_id = ?", userModel.ID, category.ID).
			Update("category_id", nil)

		if updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error removing category from tasks", updateResult.Error.Error()))
			return
		}

		c.JSON(http.StatusOK, helpers.FormatResponseWithoutData(true,
			"Default category '"+category.Name+"' removed from your tasks. "+
				fmt.Sprintf("%d tasks updated", updateResult.RowsAffected)))
		return
	}

	// Check if this category belongs to the current user
	if *category.UserID != userModel.ID {
		c.JSON(http.StatusForbidden, helpers.FormatResponseWithoutData(false, "You don't have permission to delete this category"))
		return
	}

	// This is a user-created category - check if it has associated tasks
	var taskCount int64
	initializers.DB.Model(&models.Task{}).
		Where("user_id = ? AND category_id = ?", userModel.ID, category.ID).
		Count(&taskCount)

	// Get category name before deletion
	deletedCategoryName := category.Name

	if taskCount > 0 {
		// Update all user's tasks that use this category to have null category_id
		updateResult := initializers.DB.Model(&models.Task{}).
			Where("user_id = ? AND category_id = ?", userModel.ID, category.ID).
			Update("category_id", nil)

		if updateResult.Error != nil {
			c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error updating tasks before category deletion", updateResult.Error.Error()))
			return
		}
	}

	// Perform soft delete on the user-created category
	deleteResult := initializers.DB.Delete(&category)
	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, helpers.FormatErrorResponse(false, "Error deleting category", deleteResult.Error.Error()))
		return
	}

	// Return success response with task count info
	message := "Category '" + deletedCategoryName + "' deleted successfully"
	if taskCount > 0 {
		message += fmt.Sprintf(" and removed from %d tasks", taskCount)
	}

	c.JSON(http.StatusOK, helpers.FormatResponseWithoutData(true, message))
}

func GetPendingTasks(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Find all pending tasks belonging to the current user
	var tasks []models.Task
	initializers.DB.Preload("Category").Where("user_id = ? AND is_completed = ?", userModel.ID, false).Find(&tasks)

	// Return pending tasks list
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Pending tasks retrieved successfully", helpers.FormatTasksResponse(tasks)))
}
func GetOverdueTasks(c *gin.Context) {
	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Get current time
	now := time.Now()

	// Find all overdue tasks belonging to the current user
	var tasks []models.Task
	initializers.DB.Preload("Category").
		Where("user_id = ? AND is_completed = ? AND due_date < ?", userModel.ID, false, now).
		Find(&tasks)

	// Return overdue tasks list
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Overdue tasks retrieved successfully", helpers.FormatTasksResponse(tasks)))
}
func SearchTasks(c *gin.Context) {
	// Get query parameters
	query := c.Query("q")
	category := c.Query("category")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, helpers.FormatResponseWithoutData(false, "User not authenticated"))
		return
	}

	// Convert user value to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, helpers.FormatResponseWithoutData(false, "Error retrieving user information"))
		return
	}

	// Build the query for searching tasks
	var tasks []models.Task
	dbQuery := initializers.DB.Preload("Category").Where("user_id = ?", userModel.ID)

	if query != "" {
		dbQuery = dbQuery.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if category != "" {
		dbQuery = dbQuery.Where("category_id IN (SELECT id FROM task_categories WHERE name = ? AND (user_id IS NULL OR user_id = ?))", category, userModel.ID)
	}

	dbQuery.Find(&tasks)

	// Return search results
	c.JSON(http.StatusOK, helpers.FormatResponse(true, "Search results retrieved successfully", helpers.FormatTasksResponse(tasks)))
}

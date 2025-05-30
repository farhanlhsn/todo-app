package helpers

import "todo-app/models"

type FormatSuccess struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type FormatError struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func FormatResponse(success bool, message string, data interface{}) interface{} {
	return FormatSuccess{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func FormatResponseWithoutData(success bool, message string) interface{} {
	return FormatSuccess{
		Success: success,
		Message: message,
	}
}

func FormatErrorResponse(success bool, message string, data interface{}) interface{} {
	return FormatError{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func FormatSimpleErrorResponse(success bool, message string) interface{} {
	return FormatError{
		Success: success,
		Message: message,
	}
}

type TaskResponseFormat struct {
	TaskID      uint   `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	UserID      uint   `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	DueDate     string `json:"due_date,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Category    string `json:"category,omitempty"`
}

type UserResponseFormat struct {
	ID        uint   `json:"id"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	UserSince string `json:"users_since"`
}

func FormatTaskResponse(task models.Task) TaskResponseFormat {
	var dueDateStr string
	if task.DueDate != nil {
		dueDateStr = task.DueDate.Format("2006-01-02 15:04:05")
	}

	var categoryName string
	if task.Category.ID != 0 {
		categoryName = task.Category.Name
	}

	return TaskResponseFormat{
		TaskID:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		DueDate:     dueDateStr,
		Priority:    string(task.Priority),
		Category:    categoryName,
	}
}

func FormatTasksResponse(tasks []models.Task) []TaskResponseFormat {
	formatted := make([]TaskResponseFormat, len(tasks))
	for i, task := range tasks {
		formatted[i] = FormatTaskResponse(task)
	}
	return formatted
}

func FormatUserResponse(user models.User) UserResponseFormat {
	return UserResponseFormat{
		ID:        user.ID,
		Nama:      user.Nama,
		Email:     user.Email,
		UserSince: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

type CategoriesResponseFormat struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func FormatCategoryResponse(category models.TaskCategory) CategoriesResponseFormat {
	return CategoriesResponseFormat{
		ID:   category.ID,
		Name: category.Name,
	}
}
func FormatCategoriesResponse(categories []models.TaskCategory) []CategoriesResponseFormat {
	formatted := make([]CategoriesResponseFormat, len(categories))
	for i, category := range categories {
		formatted[i] = FormatCategoryResponse(category)
	}
	return formatted
}

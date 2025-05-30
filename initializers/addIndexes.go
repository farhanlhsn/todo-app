package initializers

import (
	"log"
)

// createIndexIfNotExists safely creates an index if it doesn't exist
func createIndexIfNotExists(indexName, tableName, columns string) {
	// Check if index exists
	var count int64
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND INDEX_NAME = ?"
	DB.Raw(query, indexName).Scan(&count)

	if count == 0 {
		// Index doesn't exist, create it
		createQuery := "CREATE INDEX " + indexName + " ON " + tableName + "(" + columns + ")"
		if err := DB.Exec(createQuery).Error; err != nil {
			log.Printf("Error creating %s: %v", indexName, err)
		} else {
			log.Printf("Created index: %s", indexName)
		}
	} else {
		log.Printf("Index already exists: %s", indexName)
	}
}

// AddDatabaseIndexes creates database indexes for better performance
func AddDatabaseIndexes() {
	log.Println("Adding database indexes for performance optimization...")

	// Create indexes for tasks table
	createIndexIfNotExists("idx_tasks_user_id", "tasks", "user_id")
	createIndexIfNotExists("idx_tasks_completed", "tasks", "is_completed")
	createIndexIfNotExists("idx_tasks_due_date", "tasks", "due_date")
	createIndexIfNotExists("idx_tasks_category_id", "tasks", "category_id")
	createIndexIfNotExists("idx_tasks_priority", "tasks", "priority")

	// Create composite indexes for common queries
	createIndexIfNotExists("idx_tasks_user_completed", "tasks", "user_id, is_completed")
	createIndexIfNotExists("idx_tasks_user_due_date", "tasks", "user_id, due_date")

	// Create indexes for task_categories table
	createIndexIfNotExists("idx_task_categories_user_id", "task_categories", "user_id")
	createIndexIfNotExists("idx_task_categories_name", "task_categories", "name")

	// Create indexes for users table
	createIndexIfNotExists("idx_users_email", "users", "email")
	createIndexIfNotExists("idx_users_is_logged_in", "users", "is_logged_in")

	log.Println("Database indexes optimization completed!")
}

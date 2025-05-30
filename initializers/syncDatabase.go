package initializers

import (
	"log"
	"todo-app/models"
)

func SyncDatabase() {
	log.Println("Starting database migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.TaskCategory{},
		&models.Task{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedDefaultCategories()

	// Add database indexes for performance optimization
	AddDatabaseIndexes()

	log.Println("Database migration completed successfully")
}

func seedDefaultCategories() {
	var count int64
	DB.Model(&models.TaskCategory{}).Count(&count)

	if count == 0 {
		defaultCategories := []models.TaskCategory{
			{Name: "Work"},
			{Name: "Personal"},
			{Name: "Study"},
			{Name: "Health"},
			{Name: "Shopping"},
		}

		for _, category := range defaultCategories {
			if err := DB.Create(&category).Error; err != nil {
				log.Printf("Error creating category %s: %v", category.Name, err)
			}
		}

		log.Println("Default categories seeded successfully")
	}
}

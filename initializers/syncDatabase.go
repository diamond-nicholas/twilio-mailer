package initializers

import "github.com/nicholas/go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
package initializer

import "github.com/princesp/go-jwt/models"

func SyncDatabase() {
    DB.AutoMigrate(&models.User{})
}
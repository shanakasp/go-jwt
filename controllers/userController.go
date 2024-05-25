package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princesp/go-jwt/initializer"
	"github.com/princesp/go-jwt/models" // Import your models package
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
    // Get the email/pw from req body
    var body struct {
        Email    string
        Password string
    }
    if err := c.Bind(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return
    }

    // Hash the PW
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
        return
    }

    // Create the User
    user := models.User{Email: body.Email, Password: string(hash)}
    result := initializer.DB.Create(&user)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": result.Error.Error()})
        return
    }

    // Respond
    c.JSON(http.StatusOK, gin.H{})
}

func login(c *gin.Context) {
    // Get the email/pw from req body
    var body struct {
        Email    string
        Password string
    }
    if err := c.Bind(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return}}
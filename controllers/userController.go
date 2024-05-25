package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func Login(c *gin.Context) {
    // Get the email/pw from req body
    var body struct {
        Email    string
        Password string
    }
    if err := c.Bind(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
        return}
    
    //Look up requested user
    var user models.User
    initializer.DB.First(&user, "email = ?", body.Email)
    
    if user.ID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
        return
    }
    // Compare sent in PW saved PW hashed
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
        return
    }

    //Generate JWT token

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "sub": user.ID, 
        "exp": time.Now().Add(time.Hour * 24*30).Unix(), 
        })

    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Create Token"})
        return
    }

    // Send it back responce
    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

    c.JSON(http.StatusOK, gin.H{
       
    })
    

    }

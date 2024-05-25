package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/princesp/go-jwt/initializer"
	"github.com/princesp/go-jwt/models"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("RequireAuth")

	// Get the Cookie off req
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		if float64(time.Now().Unix()) > (*claims)["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user with the token
		var user models.User
		initializer.DB.First(&user, (*claims)["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach the user to the request context
		c.Set("user", user)

		// Continue to the next middleware/handler
		c.Next()

		fmt.Println((*claims)["foo"], (*claims)["nbf"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

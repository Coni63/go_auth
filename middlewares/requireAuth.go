package middlewares

import (
	"fmt"
	"go_auth/initializers"
	"go_auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// Check expiration time
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		c.Abort()
		return
	}

	// Convert current time to the same units as the exp claim (Unix seconds)
	now := float64(time.Now().Unix())
	if now > expFloat {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		c.Abort()
		return
	}

	// get the user ID from the token claims
	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token -- missing sub claim"})
		c.Abort()
		return
	}

	var user models.User
	initializers.DB.First(&user, "id = ?", userID)

	c.Set("user", user) // Set the user in the context for later use
	c.Next()
}

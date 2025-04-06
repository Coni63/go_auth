package middlewares

import (
	"fmt"
	"go_auth/config"
	"go_auth/initializers"
	"go_auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func cacheToken(token string, user *models.User, ttl time.Duration) {
	initializers.Cache.SetWithTTL(token, user, 1, ttl)
}

func getCachedUser(token string) (*models.User, bool) {
	val, found := initializers.Cache.Get(token)
	if !found {
		return nil, false
	}

	user, ok := val.(models.User)
	if !ok {
		return nil, false
	}
	return &user, true
}

func InvalidateUserCache(token string) {
	initializers.Cache.Del(token)
}

func GetTokenFromHeader(c *gin.Context) string {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		c.Abort()
		return ""
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	return tokenString
}

func CheckTokenValidity(tokenString string, c *gin.Context) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return nil
	}

	// Expiry check
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		c.Abort()
		return nil
	}

	if float64(time.Now().Unix()) > expFloat {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		c.Abort()
		return nil
	}

	return claims
}

func authenticateUser(tokenString string, user **models.User, c *gin.Context) *models.User {
	// Check in cache
	if val, ok := getCachedUser(tokenString); ok {
		*user = val // val is already a pointer, no need to dereference it
		return *user
	}

	claims := CheckTokenValidity(tokenString, c)
	if c.IsAborted() {
		return nil
	}

	// Load from DB
	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token -- missing sub claim"})
		c.Abort()
		return nil
	}

	// Create a new user object to populate
	*user = new(models.User)

	// Query the database with the created user pointer
	result := initializers.DB.First(*user, "id = ?", userID)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		*user = nil // Reset the pointer to nil on error
		return nil
	}

	// Cache it using exp from token
	cacheToken(tokenString, *user, config.CacheTTL)
	return *user
}

func RequireAuth(c *gin.Context) {
	tokenString := GetTokenFromHeader(c)
	if c.IsAborted() {
		return
	}

	var user *models.User
	authenticateUser(tokenString, &user, c)

	if !c.IsAborted() {
		c.Set("user", user)
		c.Next()
	}
}

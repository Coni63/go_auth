package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go_auth/config"
	"go_auth/initializers"
	"go_auth/models"
	"go_auth/serializers"
)

func GetAllUsers(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	// First, query the actual User models from the database
	var users []models.User
	if err := initializers.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Then convert to public view
	var publicUsers []serializers.UserPublicView
	for _, user := range users {
		publicUsers = append(publicUsers, user.ToPublicView())
	}

	c.JSON(http.StatusOK, publicUsers)
}

func GetUser(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	id := c.Param("id")

	// Check if it's a valid UUID
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Query the actual User model
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert to private view
	privateView := user.ToPrivateView()

	c.JSON(http.StatusOK, privateView)
}

func PutUser(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	id := c.Param("id")

	// Check if it's a valid UUID
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Query the actual User model
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert to private view
	privateView := user.ToPrivateView()

	c.JSON(http.StatusOK, privateView)
}

func PatchUser(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	id := c.Param("id")

	// Check if it's a valid UUID
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Query the actual User model
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert to private view
	privateView := user.ToPrivateView()

	c.JSON(http.StatusOK, privateView)
}

func DeleteUser(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	id := c.Param("id")

	// Check if it's a valid UUID
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Query the actual User model
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert to private view
	privateView := user.ToPrivateView()

	c.JSON(http.StatusOK, privateView)
}

func Signup(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	var body struct {
		Name     string `form:"name" binding:"required,min=3,max=50"`
		Password string `form:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindWith(&body, binding.Form); err != nil {
		// Check if it's a validation error
		if errs, ok := err.(validator.ValidationErrors); ok {
			// Create a more descriptive error message
			var errorMessages []string
			for _, e := range errs {
				errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", e.Field(), e.Tag()))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Check if the user already exists
	var existingUser models.User
	if err := initializers.DB.First(&existingUser, "name = ?", body.Name).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user := models.User{Name: body.Name, Lastname: uuid.Nil.String(), Password: string(encryptedPassword)}
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user.ToPrivateView())
}

func Login(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	var body struct {
		Name     string `form:"name" binding:"required,min=3,max=50"`
		Password string `form:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindWith(&body, binding.Form); err != nil {
		// Check if it's a validation error
		if errs, ok := err.(validator.ValidationErrors); ok {
			// Create a more descriptive error message
			var errorMessages []string
			for _, e := range errs {
				errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", e.Field(), e.Tag()))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		return
	}

	// First, find the user by name only
	var existingUser models.User
	if err := initializers.DB.First(&existingUser, "name = ?", body.Name).Error; err != nil {
		// User not found
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Now verify the password separately using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(body.Password)); err != nil {
		// Password doesn't match
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": existingUser.ID.String(),
		"exp": time.Now().Add(config.TokenTTL).Unix(), // Token expires in 24 hours
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
	if err != nil {
		log.Println("Error signing token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// If we get here, user is authenticated
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func ResetPassword(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	var body struct {
		Name     string `form:"name" binding:"required,min=3,max=50"`
		Password string `form:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindWith(&body, binding.Form); err != nil {
		// Check if it's a validation error
		if errs, ok := err.(validator.ValidationErrors); ok {
			// Create a more descriptive error message
			var errorMessages []string
			for _, e := range errs {
				errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", e.Field(), e.Tag()))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Check if the user already exists
	var existingUser models.User
	if err := initializers.DB.First(&existingUser, "name = ?", body.Name).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user := models.User{Name: body.Name, Lastname: uuid.Nil.String(), Password: string(encryptedPassword)}
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user.ToPrivateView())
}

func Validate(c *gin.Context) {
	initializers.OpsProcessed.Inc() // Increment the counter for each request

	user := c.MustGet("user").(models.User)
	c.JSON(http.StatusCreated, user.ToPrivateView())
}

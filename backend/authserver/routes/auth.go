// package routes

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
// 	"github.com/joho/godotenv"
// 	"golang.org/x/crypto/bcrypt"
// 	"main.go/config" // Use correct module path
// 	"main.go/models" // Use correct module path
// )

// var jwtSecretKey []byte

// type Claims struct {
// 	Username string `json:"username"`
// 	jwt.RegisteredClaims
// }

// func init() {
// 	// Load environment variables
// 	err := godotenv.Load("../.env")
	
// 	if err != nil {
// 		fmt.Println("there")
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Initialize JWT secret key
// 	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))
// }

// // SignupHandler handles user signup
// func SignupHandler(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Check if username already exists
// 	var existingUser models.User
// 	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
// 		return
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
// 		return
// 	}
// 	user.Password = string(hashedPassword)

// 	// Save user to the database
// 	if err := config.DB.Create(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
// }

// // LoginHandler handles user login and JWT generation
// func LoginHandler(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Retrieve user from the database
// 	var storedUser models.User
// 	if err := config.DB.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}

// 	// Compare hashed password
// 	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}

// 	// Generate JWT token
// 	expirationTime := time.Now().Add(24 * time.Hour)
// 	claims := &Claims{
// 		Username: user.Username,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtSecretKey)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
// 		return
// 	}

// 	// Respond with the JWT token
// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }

// // ProtectedHandler handles access to protected routes
// func ProtectedHandler(c *gin.Context) {
// 	// Get token from Authorization header
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
// 		return
// 	}

// 	// Parse token
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecretKey, nil
// 	})
// 	if err != nil || !token.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	// Extract claims
// 	claims, ok := token.Claims.(*Claims)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 		return
// 	}

// 	// Send protected response
// 	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %s! You have access to this protected route.", claims.Username)})
// }

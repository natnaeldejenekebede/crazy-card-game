// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gorilla/mux"
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mysql"
// 	"github.com/joho/godotenv"
// 	"golang.org/x/crypto/bcrypt"
// )

// // Card struct to represent the structure of each card
// type Card struct {
//     Suit  string `json:"suit"`
//     Value string `json:"value"`
// }

// var db *gorm.DB
// var jwtSecretKey []byte

// type User struct {
// 	ID       uint   `json:"id" gorm:"primary_key"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// type Claims struct {
// 	Username string `json:"username"`
// 	jwt.StandardClaims
// }

// func main() {
// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Read DB credentials from environment variables
// 	dbUser := os.Getenv("DBUSER")
// 	dbPassword := os.Getenv("DBPASS")
// 	dbAddr := os.Getenv("DBADDR")
// 	dbName := os.Getenv("DBNAME")

// 	// Set JWT secret key from environment variable
// 	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// 	// Create the DSN string for MySQL
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
// 		dbUser, dbPassword, dbAddr, dbName)

// 	// Connect to the database using GORM
// 	db, err = gorm.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	defer db.Close()

// 	// Migrate the schema (create table if not exists)
// 	db.AutoMigrate(&User{})

// 	r := mux.NewRouter()
// 	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
// 	r.HandleFunc("/login", LoginHandler).Methods("POST")

// 	handler := enableCORS(r) // Apply CORS middleware
// 	log.Println("Server started at :8080")
// 	http.ListenAndServe(":8080", handler)
// }
// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (* means any frontend)
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		// Handle preflight (OPTIONS) request
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func SignUpHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Check if user already exists
// 	var existingUser User
// 	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
// 		http.Error(w, "User already exists", http.StatusBadRequest)
// 		return
// 	}

// 	// Hash the password before storing it
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 		return
// 	}
// 	user.Password = string(hashedPassword)

// 	// Insert new user
// 	if err := db.Create(&user).Error; err != nil {
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode("User created successfully")
// }

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Verify user
// 	var storedUser User
// 	if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
// 		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 		return
// 	}

// 	// Compare the hashed password
// 	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
// 	if err != nil {
// 		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 		return
// 	}

// 	// Generate JWT token
// 	expirationTime := time.Now().Add(24 * time.Hour)
// 	claims := &Claims{
// 		Username: user.Username,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtSecretKey)
// 	if err != nil {
// 		http.Error(w, "Could not generate token", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send token as the response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "Login successful",
// 		"token":   tokenString,
// 	})
// }

// func authenticate(w http.ResponseWriter, r *http.Request) (*Claims, error) {
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
// 		return nil, fmt.Errorf("missing Authorization header")
// 	}

// 	tokenString := authHeader[len("Bearer "):]
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecretKey, nil
// 	})

// 	if err != nil || !token.Valid {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	claims, ok := token.Claims.(*Claims)
// 	if !ok {
// 		http.Error(w, "Could not parse token claims", http.StatusUnauthorized)
// 		return nil, fmt.Errorf("could not parse token claims")
// 	}

// 	return claims, nil
// }
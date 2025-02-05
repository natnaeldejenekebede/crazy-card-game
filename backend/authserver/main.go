package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// "strconv"
	"time"

	"math/rand"
	"sync"

	"main.go/config"
	"main.go/models"

	// "gorm.io/gorm"
	// "gorm.io/driver/mysql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	// "honnef.co/go/tools/config"
)

var roomMutex sync.Mutex

type Response struct {
	NewIP string `json:"newip"`
}
type Card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
}

type Game struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	DeckRemaining []byte `json:"deckremaining" gorm:"type:json"` // Store JSON as bytes
	TopCard       []byte `json:"topcard" gorm:"type:json"`       // Store JSON as bytes
	// StartAt       time.Time `json:"start_at"`
	// EndAt         time.Time `json:"end_at"`
	PlayersID []byte `json:"playersid" gorm:"type:json"` // Store JSON as bytes
	RoomID    uint   `json:"roomid"`
	Finished  bool   `json:"finished"`
	IP        string `json:"ip"`
	Winner    uint   `json:"winner"`
	Turn      uint   `json:"turn"`
}

type User struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	GameIDs       []byte `json:"gameids" gorm:"type:json"` // Store JSON as bytes
	CurrentGameID uint   `json:"currentgameid"`
	CurrentHand   []byte `json:"currenthand" gorm:"type:json"` // Store JSON as bytes
}

type Room struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Name          string `json:"name"`
	CurrentGameID uint   `json:"currentgameid"`
	PlayerCount   uint   `json:"playercount"`
}

// Assuming you would have the database migration function
func migrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Game{}, &Room{})
}

// var rMutex sync.Mutex
var list []uint
var db = config.InitDB()
var jwtSecretKey []byte

// Claims struct for JWT token parsing
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	migrateDatabase(db)
	fmt.Println("entered1")
	populateRooms()
	fmt.Println("entered 2")
	// Load environment variables from .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read DB credentials from environment variables
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASS")
	dbAddr := os.Getenv("DBADDR")
	dbName := os.Getenv("DBNAME")

	// Set JWT secret key from environment variable
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

	// Create the DSN string for MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbAddr, dbName)

	// Connect to the database using GORM
	db, err = gorm.Open("mysql", dsn)
	// db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	fmt.Println("entered3 ")
	// Migrate the schema (create table if not exists)
	// migrateDatabase(db)
	fmt.Println("entered4")

	r := mux.NewRouter()
	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/roomselection", RoomSelectionHandler).Methods("GET")
	r.HandleFunc("/reconnect", healthCheckAndRecovery).Methods("GET")

	handler := enableCORS(r) // Apply CORS middleware
	log.Println("Server started at :8083")
	http.ListenAndServe(":8083", handler)
	//
}

// Pre-populate the rooms with suits
func populateRooms() {
	suits := []string{"Hearts", "Diamonds", "Spades", "Clubs"}
	for i, suit := range suits {
		room := Room{
			ID:   uint(i + 1),
			Name: suit,
			// No game in the room initially
		}
		// Save the room if it doesn't exist
		if err := db.Where("id = ?", room.ID).First(&room).Error; err != nil {
			// If the room does not exist, create it
			if err := db.Create(&room).Error; err != nil {
				log.Fatalf("Error creating room: %v", err)
			}
		}
		//if err := db.FirstOrCreate(&room, Room{ID: room.ID}).Error; err != nil {
		// 	log.Fatalf("Error ensuring room exists: %v", err)
		// }

	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("eneterd 5")
	var userm models.User
	err := json.NewDecoder(r.Body).Decode(&userm)
	fmt.Println("eneterd 6", userm)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("eneterd 7")

	// Check if user already exists
	var existingUser User
	if err := db.Where("username = ?", userm.Username).First(&existingUser).Error; err == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userm.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	var user User
	user.Username = userm.Username
	user.Password = string(hashedPassword)
	gameIDsJSON, _ := json.Marshal([]uint{})
	gamehand, _ := json.Marshal([]Card{})
	user.GameIDs = gameIDsJSON
	user.CurrentGameID = uint(0)
	user.CurrentHand = gamehand
	// Insert new user
	fmt.Println(user)
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userm models.User
	err := json.NewDecoder(r.Body).Decode(&userm)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify user
	var storedUser User
	if err := db.Where("username = ?", userm.Username).First(&storedUser).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(userm.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: userm.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Send token as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func authenticate(w http.ResponseWriter, r *http.Request) (*User, *Claims, error) {
	fmt.Println("authenticate")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return nil, nil, fmt.Errorf("missing Authorization header")
	}
	fmt.Println("authenticate2")

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	fmt.Println("authenticate3")

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, nil, fmt.Errorf("invalid token")
	}
	fmt.Println("authenticate4")
	claims, ok := token.Claims.(*Claims)
	fmt.Println("authenticate", claims.Username)
	if !ok {
		http.Error(w, "Could not parse token claims", http.StatusUnauthorized)
		return nil, nil, fmt.Errorf("could not parse token claims")
	}
	fmt.Println("authenticate5")
	// Fetch the user from the database using the claims (assuming claims.ID is the user ID)
	var user User
	if err := db.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return nil, nil, fmt.Errorf("user not found")
	}
	fmt.Println("authenticate6")

	return &user, claims, nil
}
func RoomSelectionHandler(w http.ResponseWriter, r *http.Request) {
	

	roomMutex.Lock() // Lock before selecting a room
	defer roomMutex.Unlock()
	user, _, err := authenticate(w, r)
	fmt.Println("enters")
	if err != nil {
		return
	}
	suit := r.URL.Query().Get("suit")
	if suit == "" {
		http.Error(w, "Suit is required", http.StatusBadRequest)
		return
	}
	fmt.Println("enters2")
	err = db.Where("id = ? ", user.ID).First(user).Error
	if err != nil {
		return
	}
	fmt.Println("enters3")
	if user.CurrentGameID == 0 {
		fmt.Println("enter4")
		var avroom Room
		err = db.Where("name = ? ", suit).First(&avroom).Error
		if err != nil {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}
		fmt.Println("enter5")
		if avroom.PlayerCount == 3 {
			http.Error(w, "Room is full", http.StatusBadRequest)
			return
		}
		if avroom.PlayerCount == 0 {
			fmt.Println("enters6")
			playid, _ := json.Marshal([]uint{user.ID})
			var selectedServerIP string

			var gameServers []string
			err := db.Model(&Game{}).Select("DISTINCT ip").Pluck("ip", &gameServers).Error
			fmt.Println("entere7", gameServers)
			if err != nil {
				http.Error(w, "Error fetching game servers", http.StatusInternalServerError)
				return
			}

			// If no game servers exist, randomly select between two IPs
			if len(gameServers) == 0 {
				rand.Seed(time.Now().UnixNano())                              // Seed for randomness
				availableIPs := []string{"192.168.100.5:8080"}                //"192.168.100.5:8081"
				selectedServerIP = availableIPs[rand.Intn(len(availableIPs))] // Pick a random IP
			} else {
				fmt.Println("entere8")
				// Map to store unfinished game counts for each server
				serverUnfinishedGameCounts := make(map[string]int)

				// Count unfinished games for each server
				for _, ip := range gameServers {
					var unfinishedGameCount int
					err := db.Model(&Game{}).Where("ip = ? AND finished = ?", ip, false).Count(&unfinishedGameCount).Error
					if err != nil {
						http.Error(w, "Error counting unfinished games for server", http.StatusInternalServerError)
						return
					}
					serverUnfinishedGameCounts[ip] = unfinishedGameCount
				}
				fmt.Println("entere10")
				// Choose the server with the least number of unfinished games
				minUnfinishedGames := int(^uint(0) >> 1) // Max int value
				for ip, count := range serverUnfinishedGameCounts {
					if count < minUnfinishedGames {
						minUnfinishedGames = count
						selectedServerIP = ip
					}
				}

				// If no server was selected, randomly choose between the two default IPs
				if selectedServerIP == "" {
					rand.Seed(time.Now().UnixNano())
					availableIPs := []string{"192.168.100.5:8080"} //"192.168.100.5:8081"
					selectedServerIP = availableIPs[rand.Intn(len(availableIPs))]
				}
			}

			fmt.Println("enters8")
			game := &Game{
				// DeckRemaining: []byte{}, // Initialize deck, will populate this in actual game logic
				// TopCard:       []byte{},   // Initialize top card
				// StartAt: time.Now(),
				// EndAt:   time.Time{}, // Will be set when the game ends
				PlayersID: playid, // Add the user to the players list
				RoomID:    avroom.ID,
				Finished:  false,
				IP:        selectedServerIP, // Use the selected server IP
				Winner:    0,                // Will be set once the game finishes
			}
			fmt.Println("enters9")
			err = db.Create(&game).Error
			if err != nil {
				http.Error(w, "Error creating game", http.StatusInternalServerError)
				return
			}
			fmt.Println("enters11")
			user.CurrentGameID = game.ID
			err = db.Save(&user).Error
			//err = db.Model(&User{}).Where("id = ?", user.ID).Update("currentgameid", game.ID).Error
			if err != nil {
				http.Error(w, "Error updating user with current game", http.StatusInternalServerError)
				return
			}
			fmt.Println("enters`12")
			avroom.CurrentGameID = game.ID
			avroom.PlayerCount++
			err = db.Save(&avroom).Error
			if err != nil {
				http.Error(w, "Error updating room with current game", http.StatusInternalServerError)
				return
			}
			fmt.Println("enters13")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "good",
				"ip":      game.IP,
			})

		} else {
			fmt.Println("enters14")
			var game Game
			err = db.Where("id = ?", avroom.CurrentGameID).First(&game).Error
			if err != nil {
				return
			}
			// var playerIDs []uint

			// // Unmarshal existing PlayersID JSON if not empty
			// if len(game.PlayersID) > 0 {
			// 	err := json.Unmarshal(game.PlayersID, &playerIDs)
			// 	if err != nil {
			// 		fmt.Println("Error unmarshalling PlayersID:", err)
			// 		return
			// 	}
			// }

			// // Append the new user ID
			// playerIDs = append(playerIDs, user.ID)
			// fmt.Println(playerIDs)

			// // Marshal back to JSON
			// game.PlayersID, _ = json.Marshal(playerIDs)
			// game.PlayersID, _ = json.Marshal(list)
			var playerIDs []interface{}
			if len(game.PlayersID) > 0 {
				if err := json.Unmarshal(game.PlayersID, &playerIDs); err != nil {
					fmt.Println("enters15", err)
					return
				}
			}
			fmt.Println("enters14=6", playerIDs)
			// Append new player ID
			playerIDs = append(playerIDs, user.ID)
			fmt.Println("enters14=6", playerIDs)
			// Encode back to JSON
			updatedPlayersID, err := json.Marshal(playerIDs)
			if err != nil {
				fmt.Println("enters17")
				return
			}
			game.PlayersID = updatedPlayersID
			// Save back to the database
			if err := db.Save(&game).Error; err != nil {
				fmt.Println("enters18")
				fmt.Println("Error saving updated PlayersID:", err)
			}

			// err = db.Where("id = ? ", avroom.CurrentGameID).First(&game).Error
			// game.PlayersID = append(game.PlayersID, byte(user.ID))
			// fmt.Println(game.PlayersID)
			// err = db.Save(&game).Error

			avroom.PlayerCount++
			err = db.Save(&avroom).Error
			fmt.Println("enters15")
			if err != nil {
				http.Error(w, "Error updating room with current game", http.StatusInternalServerError)
				return
			}
			fmt.Println("enters16")
			user.CurrentGameID = game.ID
			err = db.Save(&user).Error
			//err = db.Model(&User{}).Where("id = ?", user.ID).Update("currentgameid", avroom.CurrentGameID).Error

			if err != nil {
				http.Error(w, "Error updating user with current game", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "good",
				"ip":      game.IP,
			})

		}

	} else {
		return
	}

}

// Health check every 30 seconds
func healthCheckAndRecovery(w http.ResponseWriter, r *http.Request) {
	// for {
	// Wait for 30 seconds before checking the server health
	// time.Sleep(30 * time.Second)

	// Fetch all game servers with unfinished games
	var gameServers []string
	err := db.Model(&Game{}).Select("DISTINCT ip").Pluck("ip", &gameServers).Error
	if err != nil {
		fmt.Println("Error fetching game servers:", err)
		// continue
	}
	if len(gameServers) <= 1 {
		gameServers = []string{
			"192.168.100.5:8080",
			"192.168.100.5:8081",
		}
	}
	// Iterate through all game servers and check their health
	for _, ip := range gameServers {
		if !isServerHealthy(ip) {

			// If the server is not responding, handle unfinished games on that server
			newip := handleUnfinishedGames(ip)
			if newip == "" {
				return
			}

			response := Response{NewIP: newip} // Example new IP
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		}
	}
}

// }

// Check if the game server is healthy by pinging it
func isServerHealthy(ip string) bool {
	resp, err := http.Get(fmt.Sprintf("http://%s/health", ip)) // Assuming the server has a /health endpoint
	if err != nil {
		fmt.Printf("Server %s is not responding\n", ip)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Server %s is healthy\n", ip)
		return true
	}

	fmt.Printf("Server %s responded with status %d\n", ip, resp.StatusCode)
	return false
}

// Handle unfinished games when the server is down
func handleUnfinishedGames(ip string) string {
	// Find unfinished games associated with this server
	var unfinishedGames []Game
	err := db.Where("ip = ? AND finished = ?", ip, false).Find(&unfinishedGames).Error
	if err != nil {
		fmt.Println("Error fetching unfinished games:", err)
		return ""
	}

	// If there are unfinished games, reassign them to a new server
	if len(unfinishedGames) > 0 {
		// Choose a new game server with the least number of unfinished games
		newIP := selectNewServerIP()
		fmt.Println(newIP)
		// Reassign the unfinished games to the new server
		for _, game := range unfinishedGames {
			game.IP = newIP
			err := db.Save(&game).Error
			if err != nil {
				fmt.Println("Error saving game with new IP:", err)
				continue
			}

		}
		return newIP

		// Notify the clients about the new IP
		// notifyClientsOfNewServer(newIP)
	}
	return ""
}

// Select a new game server IP with the least number of unfinished games
func selectNewServerIP() string {
	var gameServers []string
	err := db.Model(&Game{}).Select("DISTINCT ip").Pluck("ip", &gameServers).Error
	if err != nil {
		fmt.Println("Error fetching game servers:", err)
		return ""
	}
	if len(gameServers) <= 1 {
		gameServers = []string{
			"192.168.100.5:8080",
			"192.168.100.5:8081",
		}
	}
	// Map to store unfinished game counts for each server
	serverUnfinishedGameCounts := make(map[string]int)

	// Count unfinished games for each server
	for _, ip := range gameServers {
		var unfinishedGameCount int
		err := db.Model(&Game{}).Where("ip = ? AND finished = ?", ip, false).Count(&unfinishedGameCount).Error
		if err != nil {
			fmt.Println("Error counting unfinished games for server:", err)
			continue
		}
		serverUnfinishedGameCounts[ip] = unfinishedGameCount
	}

	// Choose the server with the least number of unfinished games
	var selectedServerIP string
	minUnfinishedGames := int(^uint(0) >> 1) // Max int value
	for ip, count := range serverUnfinishedGameCounts {
		if count < minUnfinishedGames {
			minUnfinishedGames = count
			selectedServerIP = ip
		}
	}

	// If no server is selected, pick a random default IP
	if selectedServerIP == "" {
		rand.Seed(time.Now().UnixNano())
		availableIPs := []string{"192.168.100.5:8081", "192.168.100.5:8080"}
		selectedServerIP = availableIPs[rand.Intn(len(availableIPs))]
	}

	return selectedServerIP
}

// Notify the clients about the new server IP
// func notifyClientsOfNewServer(newIP string) {
// 	// Fetch players of unfinished games
// 	var unfinishedGames []Game
// 	err := db.Where("ip = ? AND finished = ?", newIP, false).Find(&unfinishedGames).Error
// 	if err != nil {
// 		fmt.Println("Error fetching unfinished games:", err)
// 		return
// 	}

// 	// Send a notification to the players of the unfinished game (you can expand this logic based on your messaging system)
// 	for _, game := range unfinishedGames {
// 		for _, playerID := range game.PlayersID {
// 			// Encode JSON
// 			//  var new uint
// 			// playerIDnew:= json.Unmarshal(playerID,&new)
// 			// Notify the player (You may need to adjust this to send an actual message or API request)
// 			player, _ := getUserByID(uint(playerID)) // Fetch player from the database
// 			if player != nil {
// 				// Send the new IP to the player (assumes a function to notify player, e.g., via WebSocket or email)
// 				sendNotification(player, newIP)
// 			}
// 		}
// 	}
// }

// Send notification to the player (you can implement your own notification logic here)
// func sendNotification(player *User, newIP string) {

// 	// Example: Send the new server IP to the player
// 	fmt.Printf("Sending notification to player %s with new IP %s\n", player.Username, newIP)
// }

// func getUserByID(userID uint) (*User, error) {
// 	var user User
// 	// Assuming you're using GORM to query the database
// 	err := db.First(&user, userID).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("user not found: %w", err)
// 	}
// 	return &user, nil
// }

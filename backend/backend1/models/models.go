package models

import (
	// "fmt"
	// "time"

	"github.com/jinzhu/gorm"
)

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

// Migrate the user model to create the users table
func MigrateUser(db *gorm.DB) {
	// Automatically migrate the schema
	// fmt.Println("not  migrating")
	db.AutoMigrate(&User{}, &Room{}, &Game{})
	// fmt.Println("finshed migrating")
}

package config

import (
	"fmt"
	"log"
	"os"

	"main/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

// Load environment variables from the .env file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("hello")
		log.Fatal("Error loading .env file")
	}
	// fmt.Println("rrrrr")
}

// InitDB initializes the MySQL database connection using GORM
func InitDB() *gorm.DB {
	// Load environment variables first
	LoadEnv()

	// Get the MySQL credentials from the .env file
	// dbUser := "root"
	// dbPass := "yohanneshabtamu"
	// dbHost := "0.0.0.0:3306"
	// dbName := "mydb"
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbHost := os.Getenv("DBADDR")
	dbName := os.Getenv("DBNAME")

	// fmt.Println(dbHost)
	// Create the MySQL connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)

	// Connect to the MySQL database
	var errDb error
	DB, errDb = gorm.Open("mysql", dsn)
	if errDb != nil {
		log.Fatalf("Failed to connect to the database: %v", errDb)
	}
	// defer DB.Close()
	models.MigrateUser(DB)
	// Return the DB instance
	return DB
}

// CloseDB closes the database connection gracefully
func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing the database connection: %v", err)
	}
}

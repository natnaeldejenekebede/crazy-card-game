package models

import "github.com/jinzhu/gorm"

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// Migrate the user model to create the users table
func MigrateUser(db *gorm.DB) {
	// Automatically migrate the schema
	db.AutoMigrate(&User{})
}

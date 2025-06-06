package database

import (
	"fmt"
	"log"
	"os"
	"github.com/adibhar/blockchain-api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    log.Println("Initializing database...")
    err := godotenv.Load()
    if err != nil {
        log.Println("Failed to load .env file")
    }

    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    name := os.Getenv("DB_NAME")
    port := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        host, user, password, name, port)
    log.Println("Connecting to database...")

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database", err)
    }

    DB = db
    log.Println("Successfully connected to the database.")

    err = db.AutoMigrate(&models.Block{}, &models.Transaction{})
    if err != nil {
        log.Fatal("Failed to migrate database")
    }

    log.Println("Database migration completed.")
}

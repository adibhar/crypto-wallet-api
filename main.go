package main

import (
	"log"

	"github.com/adibhar/blockchain-api/database"
	"github.com/adibhar/blockchain-api/models"
	"github.com/adibhar/blockchain-api/routes"
	"github.com/adibhar/blockchain-api/services"
)

func main() {
	database.InitDB()

	var count int64
	if err := database.DB.Model(&models.Block{}).Count(&count).Error; err != nil {
		log.Fatal("Failed to check block count", err)
	}

	if count == 0 {
		genesis := services.GenerateGenesisBlock()
		if err := database.DB.Create(&genesis).Error; err != nil {
			log.Fatal("Failed to create genesis block: ", err)
		}
		log.Println("Genesis block created.")
	}

	router := routes.SetupRouter()

	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed: ", err)
	}
}

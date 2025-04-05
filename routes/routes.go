package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/adibhar/blockchain-api/database"
	"github.com/adibhar/blockchain-api/models"
)

var difficulty = 3

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/blockchain", func(c *gin.Context) {
		var blockchain []models.Block
		if err := database.DB.Preload("Transactions").Find(&blockchain).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, blockchain)
	})

	router.GET("/mine", func(c *gin.Context) {
		var lastBlock models.Block
		if err := database.DB.Last(&lastBlock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve the last block"})
			return
		}

		var transactions []models.Transaction
		if err := database.DB.Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve transactions"})
			return
		}

		newBlock := models.Block{
			Index:        lastBlock.Index + 1,
			Timestamp:    time.Now(),
			PrevHash:     lastBlock.Hash,
			Transactions: transactions,
		}

		newBlock.MineBlock(difficulty)

		if err := database.DB.Create(&newBlock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save mined block"})
			return
		}

		_ = database.DB.Delete(&transactions)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Block mined successfully!",
			"block":   newBlock,
		})
	})

	return router
}

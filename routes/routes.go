package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/adibhar/blockchain-api/database"
	"github.com/adibhar/blockchain-api/models"
	"github.com/adibhar/blockchain-api/services"
	"crypto/ecdsa"
	"math/big"
	"encoding/hex"
	"crypto/elliptic"
	"github.com/gin-contrib/cors"
	
)

const Difficulty = 3
const MiningReward = 50.0

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, //edit
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))


	router.GET("/blockchain", func(c *gin.Context) {
		var blockchain []models.Block
		if err := database.DB.Preload("Transactions").Find(&blockchain).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, blockchain)
	})

	router.GET("/balance/:address", func(c *gin.Context) {
		address := c.Param("address")

		var blockchain []models.Block
		if err := database.DB.Preload("Transactions").Find(&blockchain).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var balance float64

		for _, block := range blockchain {
			for _, tx := range block.Transactions {
				if tx.Receiver == address {
					balance += tx.Amount
				}
				if tx.Sender == address {
					balance -= tx.Amount
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"address": address,
			"balance": balance,
		})
	})

	router.POST("/transact", func(c *gin.Context) {
		var transaction models.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction data"})
			return
		}
		
		senderPubKeyBytes, err := hex.DecodeString(transaction.Sender)
		if err != nil || len(senderPubKeyBytes) != 65 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender public key"})
			return
		}
	
		x := new(big.Int).SetBytes(senderPubKeyBytes[1:33])
		y := new(big.Int).SetBytes(senderPubKeyBytes[33:65])
	
		publicKey := ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		}
	
		if !transaction.Verify(&publicKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
		var blockchain []models.Block
		if err := database.DB.Preload("Transactions").Find(&blockchain).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blockchain"})
			return
		}

		var senderBalance float64

		for _, block := range blockchain {
			for _, tx := range block.Transactions {
				if tx.Receiver == transaction.Sender {
					senderBalance += tx.Amount
				}
				if tx.Sender == transaction.Sender {
					senderBalance -= tx.Amount
				}
			}
		}
	
		// senderBalance := models.CalculateBalance(transaction.Sender, blockchain)
		if senderBalance < transaction.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
			return
		}
	
		err = database.DB.Create(&transaction).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add transaction",
				"details": err.Error(),
			})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction created successfully",
			"tx":      transaction,
		})

	})
	router.POST("/mine/:address", func(c *gin.Context) {
		address := c.Param("address");
		
		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Need miner address"})
			return;
		}

		var lastBlock models.Block;	
		if err := database.DB.Last(&lastBlock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve last block"})
			return;
		}

		var transactions []models.Transaction
		if err := database.DB.Where("block_id IS NULL").Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve transactions"})
			return
		}

		rewardTx := models.Transaction{
			Sender:   "SYSTEM",
			Receiver: address,
			Amount:   MiningReward,
			Signature: "reward",
		}
		transactions = append(transactions, rewardTx)

		newBlock := services.GenerateBlock(lastBlock, transactions);
		if err := database.DB.Create(&newBlock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid new block saved"})
			return;
		} 

		c.JSON(http.StatusCreated, gin.H{
			"Message": "Block mined successfully, reward granted",
			"block" : newBlock,
		})


	})

	//! DEV ROUTE
	router.POST("/reset", func(c *gin.Context) {
		if err := database.DB.Exec("DELETE FROM transactions").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete transactions", "details": err.Error()})
			return
		}
		if err := database.DB.Exec("DELETE FROM blocks").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete blocks", "details": err.Error()})
			return
		}

		if err := database.DB.Exec("ALTER SEQUENCE transactions_id_seq RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset transaction ID sequence", "details": err.Error()})
			return
		}
		if err := database.DB.Exec("ALTER SEQUENCE blocks_id_seq RESTART WITH 1").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset block ID sequence", "details": err.Error()})
			return
		}
	
		genesis := services.GenerateGenesisBlock()
		if err := database.DB.Create(&genesis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create genesis block", "details": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"message":       "Blockchain reset with genesis block",
			"genesis_block": genesis,
		})
	})
	
	//! END DEV ROUTE 


	router.GET("/mine", func(c *gin.Context) {
		var lastBlock models.Block
		if err := database.DB.Last(&lastBlock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve the last block"})
			return
		}

		var transactions []models.Transaction
		if err := database.DB.Where("block_id IS NULL").Find(&transactions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve transactions"})
			return
		}

		if len(transactions) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No transactions to mine"})
			return
		}

		newBlock := services.GenerateBlock(lastBlock, transactions)
		if err := database.DB.Create(&newBlock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save mined block"})
			return
		}

		for _, tx := range transactions {
			tx.BlockID = &newBlock.ID
			if err := database.DB.Save(&tx).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction with block ID"})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Block mined successfully!",
			"block":   newBlock,
		})
	})

	router.GET("/pending", func(c *gin.Context) {
		var pending []models.Transaction
	
		if err := database.DB.Where("block_id IS NULL").Find(&pending).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve pending transactions"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"pending_transactions": pending,
		})
	})

	router.GET("/verify", func(c *gin.Context) {
		var blockchain []models.Block
		if err := database.DB.Preload("Transactions").Find(&blockchain).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		for i := 1; i < len(blockchain); i++ {
			prevBlock := blockchain[i-1]
			currBlock := blockchain[i]
	
			if !services.IsBlockValid(currBlock, prevBlock) {
				c.JSON(http.StatusOK, gin.H{
					"valid": false,
					"error": "Block validation failed",
					"index": i,
				})
				return
			}
		}
	
		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"message": "Blockchain is valid",
		})
	})
	
	return router
}

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
	
)

const Difficulty = 3
const MiningReward = 50.0

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

	return router
}

package services

import (
	"time"
	"github.com/adibhar/blockchain-api/models"
)

const Difficulty = 3

func GenerateBlock(oldBlock models.Block, transactions []models.Transaction) models.Block {
	timestamp := time.Now()

	newBlock := models.Block{
		Index:        oldBlock.Index + 1,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     oldBlock.Hash,
	}

	newBlock.MineBlock(Difficulty)

	return newBlock
}

func GenerateGenesisBlock() models.Block {
	genesisBlock := models.Block{
		Index:     0,
		Timestamp: time.Now(),
		PrevHash:  "NONE",
		Hash:      "0000",
		Nonce:     0,
	}

	return genesisBlock
}

func IsBlockValid(newBlock, oldBlock models.Block) bool {
	if oldBlock.Index+1 != newBlock.Index || oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	txData := newBlock.SerializeTransactions(newBlock.Transactions)
	calculatedHash, _ := models.ProofOfWork(newBlock.Index, newBlock.PrevHash, newBlock.Timestamp.String(), txData, Difficulty)

	return calculatedHash == newBlock.Hash
}


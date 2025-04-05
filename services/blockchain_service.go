package services

import (
	"time"
	"strings"
	"fmt"
	"github.com/adibhar/blockchain-api/models"
)

const Difficulty = 3

// GenerateBlock creates a new block using previous block data and list of transactions
func GenerateBlock(oldBlock models.Block, transactions []models.Transaction) models.Block {
	timestamp := time.Now()

	// Convert transactions to string for hashing
	txData := serializeTransactions(transactions)

	hash, nonce := models.ProofOfWork(
		oldBlock.Index+1,
		oldBlock.Hash,
		timestamp.String(),
		txData,
		Difficulty,
	)

	newBlock := models.Block{
		Index:        oldBlock.Index + 1,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     oldBlock.Hash,
		Hash:         hash,
		Nonce:        nonce,
	}

	return newBlock
}

// GenerateGenesisBlock returns the initial block in the blockchain
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

// IsBlockValid validates a new block against the previous block
func IsBlockValid(newBlock, oldBlock models.Block) bool {
	if oldBlock.Index+1 != newBlock.Index || oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	txData := serializeTransactions(newBlock.Transactions)

	calculatedHash, _ := models.ProofOfWork(
		newBlock.Index,
		newBlock.PrevHash,
		newBlock.Timestamp.String(),
		txData,
		Difficulty,
	)

	return calculatedHash == newBlock.Hash
}

// serializeTransactions converts transactions to a string for hashing
func serializeTransactions(transactions []models.Transaction) string {
	var sb strings.Builder
	for _, tx := range transactions {
		sb.WriteString(tx.Sender)
		sb.WriteString(tx.Receiver)
		sb.WriteString(fmt.Sprintf("%f", tx.Amount))
	}
	return sb.String()
}

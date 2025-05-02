package services

import (
	"time"
	"github.com/adibhar/blockchain-api/models"
	//"fmt"
	//needed for debugging 
)

const Difficulty = 3

func GenerateBlock(oldBlock models.Block, transactions []models.Transaction) models.Block {
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	

	newBlock := models.Block{
		Index:        oldBlock.Index + 1,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     oldBlock.Hash,
	}
	newBlock.MineBlock(Difficulty, timestamp)

	return newBlock
}

func GenerateGenesisBlock() models.Block {
	genesisBlock := models.Block{
		Index:     0,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
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

	//debugging stuff
	timestampStr := newBlock.Timestamp
	txData := newBlock.SerializeTransactions(newBlock.Transactions)
	// fmt.Println("hashmatching DEBUGGING")
	// fmt.Println("Index:        ", newBlock.Index)
	// fmt.Println("PrevHash:     ", newBlock.PrevHash)
	// fmt.Println("Timestamp:    ",  timestampStr)
	// fmt.Println("TxData:       ", txData)
	// fmt.Println("Stored Hash:  ", newBlock.Hash)
	calculatedHash, _ := models.ProofOfWork(newBlock.Index, newBlock.PrevHash, timestampStr, txData, Difficulty)

	return calculatedHash == newBlock.Hash
}


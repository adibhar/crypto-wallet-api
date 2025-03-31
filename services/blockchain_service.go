package services

import (
	"github.com/adibhar/blockchain-api/models"
	// "blockchain-api/db" 
	//TODO: Uncomment above line after finishing database
	"time"
)

const Difficulty = 3;

func GenerateBlock(oldBlock models.Block, transactions string) models.Block {
	timestamp := time.Now()

	hash, nonce := models.ProofOfWork(oldBlock.Index + 1, transactions, oldBlock.Hash, timestamp.String(), Difficulty);

	newBlock := models.Block{
		Index:        oldBlock.Index + 1,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     oldBlock.Hash,
		Hash:         hash,
		Nonce:        nonce,
	}
	
	return newBlock;
}

//TODO: RESTRUCTURE CODE BELOW

var Blockchain []models.Block

func GenerateGenesisBlock() models.Block {
	genesisBlock := models.Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: "Genesis Block",
		PrevHash:     "NONE",
		Hash:         "0000",
		Nonce:        0,
	}
	return genesisBlock
}


func IsBlockValid(newBlock, oldBlock models.Block) bool {
	if oldBlock.Index + 1 != newBlock.Index || oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	calculatedHash, _ := models.ProofOfWork(newBlock.Index, newBlock.Transactions, newBlock.PrevHash, newBlock.Timestamp.String(), Difficulty)
	return calculatedHash == newBlock.Hash
}



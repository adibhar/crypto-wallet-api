package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	ID            uint           `gorm:"primaryKey"`
	Index         int            `json:"index"`
	Timestamp     time.Time      `json:"timestamp"`
	PrevHash      string         `json:"prev_hash"`
	Hash          string         `json:"hash"`
	Nonce         int            `json:"nonce"`
	Transactions  []Transaction `gorm:"foreignKey:BlockID"`
}

func (b *Block) MineBlock(difficulty int) {
	txData := b.SerializeTransactions(b.Transactions)
	hash, nonce := ProofOfWork(b.Index, b.PrevHash, b.Timestamp.String(), txData, difficulty)
	b.Hash = hash
	b.Nonce = nonce
}

func ProofOfWork(index int, prevHash, timestamp, txData string, difficulty int) (string, int) {
	nonce := 0
	var hash string

	for {
		record := fmt.Sprintf("%d%s%s%s%d", index, timestamp, prevHash, txData, nonce)
		hash = calculateHash(record)

		if hash[:difficulty] == strings.Repeat("0", difficulty) {
			break
		}
		nonce++
	}

	return hash, nonce
}

func calculateHash(data string) string {
	rawHash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", rawHash[:])
}

func (b *Block) SerializeTransactions(transactions []Transaction) string {
	var sb strings.Builder
	for _, tx := range transactions {
		sb.WriteString(tx.Sender)
		sb.WriteString(tx.Receiver)
		sb.WriteString(fmt.Sprintf("%f", tx.Amount))
	}
	return sb.String()
}

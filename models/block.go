package models;

import (
	"time"
)

type Block struct {
	Index         int         `json:"index"`
	Timestamp     time.Time   `json:"timestamp"`
	Transactions  string      `json:"transactions"`
	PrevHash      string      `json:"prev_hash"`
	Hash          string      `json:"hash"`
	Nonce         int         `json:"nonce"`
}
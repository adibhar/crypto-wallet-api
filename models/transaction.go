package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
	"math/big"
	"encoding/hex"
)

type Transaction struct {
	ID            uint      `gorm:"primaryKey"`
	Sender       string    `json:"sender"`
	Receiver     string    `json:"receiver"`
	Amount       float64   `json:"amount"`
	Signature    string    `json:"signature"`
}

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func (t* Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	data := t.Sender + t.Receiver + fmt.Sprintf("%f", t.Amount);
	hash := sha256.Sum256([]byte(data));
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return err;
	}
	t.Signature = fmt.Sprintf("%x:%x", r, s)
	return nil;
}

func (t *Transaction) Verify(publicKey *ecdsa.PublicKey) bool {
	data := t.Sender + t.Receiver + fmt.Sprintf("%f", t.Amount)
	hash := sha256.Sum256([]byte(data))

	parts := strings.Split(t.Signature, ":")
	if len(parts) != 2 {
		return false
	}
	rBytes, err1 := hex.DecodeString(parts[0])
	sBytes, err2 := hex.DecodeString(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}

	var r, s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)

	return ecdsa.Verify(publicKey, hash[:], &r, &s)
}

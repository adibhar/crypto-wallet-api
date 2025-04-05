package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

type Transaction struct {
	ID        uint   `gorm:"primaryKey"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature"`
	BlockID   uint   `json:"block_id"`
	Block     Block  `gorm:"foreignKey:BlockID"`
}

func (t *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {

	amountInCents := int64(t.Amount * 100)
	data := t.Sender + t.Receiver + fmt.Sprintf("%d", amountInCents)


	hash := sha256.Sum256([]byte(data))

	
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return err
	}


	t.Signature = fmt.Sprintf("%x:%x", r, s)
	return nil
}


func (t *Transaction) Verify(publicKey *ecdsa.PublicKey) bool {

	amountInCents := int64(t.Amount * 100)
	data := t.Sender + t.Receiver + fmt.Sprintf("%d", amountInCents)

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

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}


func ExtractPublicKeyFromPrivate(privateKey *ecdsa.PrivateKey) *ecdsa.PublicKey {
	return &privateKey.PublicKey
}

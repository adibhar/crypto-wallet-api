package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Transaction struct {
	Sender       string    `json:"sender"`
	Receiver     string    `json:"reciever"`
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

	var r, s big.Int
	fmt.Sscanf(t.Signature, "%x:%x", &r, &s);
	return ecdsa.Verify(publicKey, hash[:], &r, &s);
}
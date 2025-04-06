package main

import (
	"fmt"
	"log"
	"github.com/adibhar/blockchain-api/models"
	"crypto/elliptic"
)

func main() {

	senderPrivateKey, err := models.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Error generating sender keys: %v", err)
	}
	senderPublicKey := models.ExtractPublicKeyFromPrivate(senderPrivateKey)

	receiverPrivateKey, err := models.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Error generating receiver keys: %v", err)
	}
	receiverPublicKey := models.ExtractPublicKeyFromPrivate(receiverPrivateKey)

	senderPubKeyBytes := elliptic.Marshal(elliptic.P256(), senderPublicKey.X, senderPublicKey.Y)
	receiverPubKeyBytes := elliptic.Marshal(elliptic.P256(), receiverPublicKey.X, receiverPublicKey.Y)


	senderPublicKeyStr := fmt.Sprintf("%x", senderPubKeyBytes)
	receiverPublicKeyStr := fmt.Sprintf("%x", receiverPubKeyBytes)

	fmt.Println("Sender Public Key:", senderPublicKeyStr)
	fmt.Println("Receiver Public Key:", receiverPublicKeyStr)


	tx := &models.Transaction{
		Sender:   senderPublicKeyStr,
		Receiver: receiverPublicKeyStr, 
		Amount:   10.0,
	}

	err = tx.Sign(senderPrivateKey)
	if err != nil {
		log.Fatalf("Error signing transaction: %v", err)
	}


	fmt.Println("\nTransaction Data (for /transact API):")
	fmt.Printf("Sender: %s\n", tx.Sender)
	fmt.Printf("Receiver: %s\n", tx.Receiver)
	fmt.Printf("Amount: %.2f\n", tx.Amount)
	fmt.Printf("Signature: %s\n", tx.Signature)

	//json body format
	// {
	//   "sender": "<sender public key>",
	//   "receiver": "<receiver public key>",
	//   "amount": xx.x,
	//   "signature": "<hex signature>"
	// }
}

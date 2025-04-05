package main

import (
	"log"

	"github.com/adibhar/blockchain-api/database"
	"github.com/adibhar/blockchain-api/routes"
)

func main() {
	database.InitDB();
	routes := routes.SetupRouter();
	//uncomment above line after finishing routes
	
	log.Fatal(routes.Run(":8080"))




	//TEST CODE FOR BLOCKCHAIN: DELETE

	// genesisBlock := services.GenerateGenesisBlock()
	// services.Blockchain = append(services.Blockchain, genesisBlock)
	// fmt.Println("Genesis Block Created:", genesisBlock)

	// for i := 1; i <= 3; i++ {
	// 	newBlock := services.GenerateBlock(services.Blockchain[len(services.Blockchain)-1], fmt.Sprintf("filler %d", i))

	// 	if services.IsBlockValid(newBlock, services.Blockchain[len(services.Blockchain)-1]) {
	// 		services.Blockchain = append(services.Blockchain, newBlock)
	// 		fmt.Printf("Block %d added: %+v\n", newBlock.Index, newBlock)
	// 	} else {
	// 		fmt.Println("Invalid block.")
	// 	}

	// 	time.Sleep(1 * time.Second)
	// }

	// fmt.Println("\nBlockchain:")
	// for _, block := range services.Blockchain {
	// 	fmt.Printf("Index: %d, Hash: %s, PrevHash: %s, Transactions: %s\n", block.Index, block.Hash, block.PrevHash, block.Transactions)
	// }

	//TODO: Test transactions and signatures
}

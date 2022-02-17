package main

import (
	"fmt"
	"time"
)

type Block struct {
	nonce        int
	transactions []string
	previousHash string
	timestamp    int64
}

func newBlock(nonce int, previousHash string) *Block {
	b := Block{
		nonce:        nonce,
		previousHash: previousHash,
		transactions: []string{},
		timestamp:    time.Now().UnixNano(),
	}

	return &b
}

func Print(b *Block) {
	fmt.Printf("nonce: %v\n", b.nonce)
	fmt.Printf("transactions: %v\n", b.transactions)
	fmt.Printf("previousHash: %v\n", b.previousHash)
	fmt.Printf("timestamp: %v\n", b.timestamp)
}

func main() {
	b := newBlock(0, "init hash")
	Print(b)
}

package main

import (
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	transactions []string
	previousHash string
	timestamp    int64
}

func NewBlock(nonce int, previousHash string) *Block {
	b := Block{
		nonce:        nonce,
		previousHash: previousHash,
		transactions: []string{},
		timestamp:    time.Now().UnixNano(),
	}

	return &b
}

func (b *Block) Print() {
	fmt.Printf("nonce: %v\n", b.nonce)
	fmt.Printf("transactions: %v\n", b.transactions)
	fmt.Printf("previousHash: %v\n", b.previousHash)
	fmt.Printf("timestamp: %v\n", b.timestamp)
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	// bc := Blockchain{}
	bc.CreateBlock(0, "Init Hash")

	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)

	return b
}

func (bc *Blockchain) Print() {
	for i, Block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		Block.Print()
	}
	fmt.Printf("%s \n", strings.Repeat("* ", 25))
}

func main() {
	// b := NewBlock(0, "init hash")
	// b.Print()

	bc := NewBlockchain()
	// fmt.Println(bc)
	bc.Print()

	bc.CreateBlock(5, "hash 1")
	bc.Print()

	bc.CreateBlock(2, "hash 2")
	bc.Print()
}

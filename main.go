package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	b := Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: []string{},
	}

	return &b
}

// func NewBlock(nonce int, previousHash [32]byte) *Block {
// 	b := new(Block)
// 	b.timestamp = time.Now().UnixNano()
// 	b.nonce = nonce
// 	b.previousHash = previousHash
// 	return b
// }

func (b *Block) Print() {
	fmt.Printf("timestamp: %v\n", b.timestamp)
	fmt.Printf("nonce: %v\n", b.nonce)
	fmt.Printf("previousHash: %x\n", b.previousHash)
	fmt.Printf("transactions: %v\n", b.transactions)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previousHash"`
		Transactions []string `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	// fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	// bc := Blockchain{}
	bc.CreateBlock(0, b.Hash())

	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
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
	bc := NewBlockchain()
	bc.Print()

	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(5, previousHash)
	bc.Print()

	previousHash = bc.LastBlock().Hash()
	bc.CreateBlock(2, previousHash)
	bc.Print()
}

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const MINING_DIFFICULTY = 3
const MINING_SENDER = "THE BLOCKCHAIN"
const MINING_REWARD = 1.0

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

// func NewBlock(nonce int, previousHash [32]byte) *Block {
// 	b := Block{
// 		timestamp:    time.Now().UnixNano(),
// 		nonce:        nonce,
// 		previousHash: previousHash,
// 		transactions: []string{},
// 	}

// 	return &b
// }

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp: %v\n", b.timestamp)
	fmt.Printf("nonce: %v\n", b.nonce)
	fmt.Printf("previousHash: %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

// struct의 각 필드 마지막에 json으로 변경했을 때 필드 이름을 지정할 수 있다.
// struct의 각 필드는 대문자로 시작해서 public 필드로 만들어주어야 접근이 가능하다.
// `json:"timestamp"` 이 없으면 marshal 결과 필드이름은 Timestamp 이다.
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previousHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b) // func (b *Block) MarshalJSON() 호출
	// fmt.Println("!!!")
	// fmt.Println(m)
	// fmt.Println(sha256.Sum256([]byte(m)))
	// fmt.Println("")

	// 특정 Block을 json으로 변경하려고 했는데, {} 가 출력됨
	// Block 구조체의 필드들이 private 이기 때문(lowercase)
	// https://stackoverflow.com/questions/26327391/json-marshalstruct-returns
	// https://stackoverflow.com/questions/21825322/why-golang-cannot-generate-json-from-struct-with-front-lowercase-character

	// 따라서 별도의 marshal()을 정의: MarshalJSON

	// fmt.Println(string(m))
	return sha256.Sum256([]byte(m)) // sha256.Sum256()는 해시 값을 만드는 함수
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	address         string
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func NewBlockchain(address string) *Blockchain {
	// b := &Block{}
	b := new(Block)
	bc := new(Blockchain)
	// bc := Blockchain{}
	bc.address = address
	bc.CreateBlock(0, b.Hash())

	return bc
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := newTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, newTransaction(t.senderAddress, t.recipientAddress, t.value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		timestamp:    0,
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	// fmt.Println(guessHashStr)

	return zeros == guessHashStr[:difficulty]
}

func (bc *Blockchain) ProofOfWork() int {
	nonce := 0
	lastBlock := bc.LastBlock()
	for !bc.ValidProof(nonce, lastBlock.previousHash, bc.CopyTransactionPool(), MINING_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.address, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()

	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")

	return true
}

func (bc *Blockchain) CalculateTotalAmount(address string) float32 {
	var totalAmount float32 = 0.0

	for _, b := range bc.chain {
		for _, t := range b.transactions {
			if address == t.recipientAddress {
				totalAmount += t.value
			} else if address == t.senderAddress {
				totalAmount -= t.value
			}
		}
	}

	return float32(totalAmount)
}

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

func newTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{senderAddress: sender, recipientAddress: recipient, value: value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" Sender Address:    %v\n", t.senderAddress)
	fmt.Printf(" Recipient Address: %v\n", t.recipientAddress)
	fmt.Printf(" Value:             %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"senderAddress"`
		RecipientAddress string  `json:"recipientAddress"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Value:            t.value,
	})
}

func (bc *Blockchain) Print() {
	for i, Block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		Block.Print()
	}
	fmt.Printf("%s \n", strings.Repeat("* ", 25))
}

func main() {
	myAddress := "my_blockchain_address"
	bc := NewBlockchain(myAddress)
	bc.Print()

	bc.AddTransaction("A", "B", 1.0)
	// previousHash := bc.LastBlock().Hash()
	// nonce := bc.ProofOfWork()
	// bc.CreateBlock(nonce, previousHash)
	bc.Mining()
	bc.Print()

	bc.AddTransaction("C", "D", 2.0)
	bc.AddTransaction("X", "Y", 3.0)
	// previousHash = bc.LastBlock().Hash()
	// nonce = bc.ProofOfWork()
	// bc.CreateBlock(nonce, previousHash)
	bc.Mining()
	bc.Print()

	fmt.Println()
	fmt.Printf("Mine: %.1f\n", bc.CalculateTotalAmount(myAddress))
	fmt.Printf("C   : %.1f\n", bc.CalculateTotalAmount("C"))
	fmt.Printf("D   : %.1f\n", bc.CalculateTotalAmount("D"))
}

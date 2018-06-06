package block

import (
	"time"
)

// Block type
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates a new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	work := NewProofOfWork(block)
	nonce, hash := work.run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

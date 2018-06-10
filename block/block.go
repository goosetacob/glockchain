package block

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/goosetacob/glockchain/transaction"
)

// Block type
type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates a new block
func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	work := NewProofOfWork(block)
	nonce, hash := work.run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// HashTransactions builds a hash of all the trasnsactions in this block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

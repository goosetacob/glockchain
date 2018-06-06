package blockchain

import "github.com/goosetaco/glockchain/block"

// Blockchain data structure
type Blockchain struct {
	blocks []*block.Block
}

// AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := block.NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// GetBlocks returns slice of blocks
func (bc *Blockchain) GetBlocks() []*block.Block {
	return bc.blocks
}

// NewGenesisBlock creates a genesis block
func newGenesisBlock() *block.Block {
	return block.NewBlock("Genesis Block", []byte{})
}

// NewBlockchain initializes a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{newGenesisBlock()}}
}

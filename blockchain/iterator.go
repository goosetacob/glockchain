package blockchain

import (
	"github.com/goosetacob/glockchain/block"
	"github.com/goosetacob/glockchain/datastore"
	"github.com/sirupsen/logrus"
)

// Iterator to inspect persitent storage
type Iterator struct {
	currentHash []byte
	database    *datastore.DataStore
}

// NewIterator build a new iteraratoe for a Blockchain
func (chain *Blockchain) NewIterator() *Iterator {
	return &Iterator{chain.leadingHash, chain.datastore}
}

// Next iterates onto the next hash in the chain
func (itr *Iterator) Next() *block.Block {
	var nextBlock *block.Block

	// get block given hash
	nextBlock, err := itr.database.GetBlock(itr.currentHash)
	if err != nil {
		logrus.Errorf("cannot read form database", err)
	}

	// update currentHash
	itr.currentHash = nextBlock.PrevBlockHash

	return nextBlock
}

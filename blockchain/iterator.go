package blockchain

import (
	"github.com/boltdb/bolt"
	"github.com/goosetacob/glockchain/block"
	"github.com/sirupsen/logrus"
)

// Iterator to inspect persitent storage
type Iterator struct {
	currentHash []byte
	database    *bolt.DB
}

// NewIterator build a new iteraratoe for a Blockchain
func (chain *Blockchain) NewIterator() *Iterator {
	return &Iterator{chain.leadingHash, chain.database}
}

// Next iteratrs onto the next hash in the chain
func (itr *Iterator) Next() *block.Block {
	var nextBlock *block.Block

	// get the block currentHash refers too
	if err := itr.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(itr.currentHash)
		nextBlock = deserializeBlock(encodedBlock)
		return nil
	}); err != nil {
		logrus.Errorf("cannot read form database", err)
	}

	// update currentHash
	itr.currentHash = nextBlock.PrevBlockHash

	return nextBlock
}

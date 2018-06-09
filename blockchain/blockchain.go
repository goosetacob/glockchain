package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/goosetacob/glockchain/block"
	"github.com/sirupsen/logrus"
)

// dbfile is the diskfile we're going to store data in
const dbfile = "glockchain_ledger"

// blocksBucket is the BoltDB bucket that's going to store block-specific data
const blocksBucket = "blocks"

// Blockchain data structure
type Blockchain struct {
	leadingHash []byte
	database    *bolt.DB
}

// NewBlockchain initializes a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	var newLeadHash []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		logrus.Errorf("cannot open dbfile: %v\n", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := newGenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, serializeBlock(genesis))
			err = b.Put([]byte("l"), genesis.Hash)
			newLeadHash = genesis.Hash
		} else {
			newLeadHash = b.Get([]byte("l"))
		}

		return nil
	})

	return &Blockchain{newLeadHash, db}
}

// Shutdown saves and shutdsdown blockchain program
func (chain *Blockchain) Shutdown() {
	chain.database.Close()
}

// AddBlock adds a block to the blockchain
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// read the last hash
	if err := chain.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	}); err != nil {
		logrus.Errorf("can't view database: %v", err)
	}

	// create new lbock with data
	newBlock := block.NewBlock(data, lastHash)

	// send update to db with new block appended
	if err := chain.database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		// put new block in db
		if err := bucket.Put(newBlock.Hash, serializeBlock(newBlock)); err != nil {
			return err
		}

		// update last hash
		if err := bucket.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}

		// update chain lead
		chain.leadingHash = newBlock.Hash
		return nil
	}); err != nil {
		logrus.Errorf("can't update database: %v", err)
	}
}

// NewGenesisBlock creates a genesis block
func newGenesisBlock() *block.Block {
	return block.NewBlock("Genesis Block", []byte{})
}

// serializeBLock converts a block struct into a byte array
func serializeBlock(b *block.Block) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		logrus.Errorf("error encoding: %v", err)
	}

	return result.Bytes()
}

// deserializeBlock converts a byte array into a block struct
func deserializeBlock(d []byte) *block.Block {
	var block block.Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		logrus.Errorf("error decoding: %v", err)
	}

	return &block
}

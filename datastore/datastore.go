package datastore

import (
	"bytes"
	"encoding/gob"
	"fmt"

	bolt "github.com/coreos/bbolt"
	"github.com/goosetacob/glockchain/block"
)

// DataStore struct to store data to disk
type DataStore struct {
	db *bolt.DB
}

// dbfile is the diskfile we're going to store data in
const dbfile = "glockchain.db"

// blocksBucket is the BoltDB bucket that's going to store block-specific data
const blocksBucket = "blocks"

// GetDataStore gets a connection to storage
func GetDataStore() (*DataStore, error) {
	boldDB, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot open dbfile: %v", err)
	}

	return &DataStore{boldDB}, nil
}

// Close down the datastore connection
func (datastore *DataStore) Close() {
	datastore.db.Close()
}

// GetLastHash gets the lastHash
func (datastore *DataStore) GetLastHash() ([]byte, error) {
	var lastHash []byte
	if err := datastore.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			return fmt.Errorf("bucket does not exist")
		}
		lastHash = bucket.Get([]byte("l"))
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error getting last hash %v", err)
	}
	return lastHash, nil
}

// GetBlock given hash
func (datastore *DataStore) GetBlock(blockHash []byte) (*block.Block, error) {
	var nextBlock *block.Block
	// get the block currentHash refers too
	if err := datastore.db.View(func(tx *bolt.Tx) error {
		var err error
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(blockHash)
		nextBlock, err = deserializeBlock(encodedBlock)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("cannot read form database %v", err)
	}
	return nextBlock, nil
}

// AddBlock stores a new block
func (datastore *DataStore) AddBlock(newBlock *block.Block) ([]byte, error) {
	var newLeadHash []byte
	if err := datastore.db.Update(func(tx *bolt.Tx) error {
		var err error
		var serializedBlock []byte

		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return fmt.Errorf("error creating bucket %v", err)
			}
		}

		if serializedBlock, err = serializeBlock(newBlock); err != nil {
			return fmt.Errorf("error serializing block %v", err)
		}
		if err := bucket.Put(newBlock.Hash, serializedBlock); err != nil {
			return fmt.Errorf("error PUT db %v", err)
		}
		if err := bucket.Put([]byte("l"), newBlock.Hash); err != nil {
			return fmt.Errorf("error PUT db %v", err)
		}

		newLeadHash = newBlock.Hash
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error storing new block %v", err)
	}
	return newLeadHash, nil
}

// serializeBLock converts a block struct into a byte array
func serializeBlock(b *block.Block) ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("error encoding: %v", err)
	}

	return result.Bytes(), nil
}

// deserializeBlock converts a byte array into a block struct
func deserializeBlock(d []byte) (*block.Block, error) {
	var block block.Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %v", err)
	}

	return &block, nil
}

package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"

	"github.com/boltdb/bolt"
	"github.com/goosetacob/glockchain/block"
	"github.com/goosetacob/glockchain/transaction"
	"github.com/sirupsen/logrus"
)

// dbfile is the diskfile we're going to store data in
const dbfile = "glockchain_ledger"

// blocksBucket is the BoltDB bucket that's going to store block-specific data
const blocksBucket = "blocks"

// genesisCoinbaseData
const genesisCoinbaseData = "GooseCoin will ICO before NandoCoin"

// Blockchain data structure
type Blockchain struct {
	leadingHash []byte
	database    *bolt.DB
}

// NewBlockchain initializes a new blockchain with a genesis block
func NewBlockchain(address string) *Blockchain {
	var newLeadHash []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		logrus.Errorf("cannot open dbfile: %v\n", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			coinbaseTransaction := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := newGenesisBlock(coinbaseTransaction)
			bucket, _ := tx.CreateBucket([]byte(blocksBucket))
			err = bucket.Put(genesis.Hash, serializeBlock(genesis))
			err = bucket.Put([]byte("l"), genesis.Hash)
			newLeadHash = genesis.Hash
		} else {
			newLeadHash = bucket.Get([]byte("l"))
		}
		return nil
	})

	return &Blockchain{newLeadHash, db}
}

// Shutdown saves and shutdsdown blockchain program
func (chain *Blockchain) Shutdown() {
	chain.database.Close()
}

// MineBlock mines a new block with the provided transactions
func (chain *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	var lastHash []byte

	err := chain.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		logrus.Panic(err)
	}

	newBlock := block.NewBlock(transactions, lastHash)

	err = chain.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, serializeBlock(newBlock))
		if err != nil {
			logrus.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			logrus.Panic(err)
		}

		chain.leadingHash = newBlock.Hash
		return nil
	})
}

// FindUTXO gets the outputs of unspend transactions
func (chain *Blockchain) FindUTXO(address string) []transaction.TXOutput {
	var utxo []transaction.TXOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				utxo = append(utxo, out)
			}
		}
	}
	return utxo
}

// FindUnspentTransactions looks for upspent money in Glockchain
func (chain *Blockchain) FindUnspentTransactions(address string) []transaction.Transaction {
	var unspentTXs []transaction.Transaction
	spentTXOs := make(map[string][]int)
	chainItr := chain.NewIterator()

	for {
		block := chainItr.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindSpendableOutputs does something
func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// NewGenesisBlock creates a genesis block
func newGenesisBlock(coinbase *transaction.Transaction) *block.Block {
	return block.NewBlock([]*transaction.Transaction{coinbase}, []byte{})
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

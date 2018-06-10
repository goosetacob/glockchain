package blockchain

import (
	"encoding/hex"

	"github.com/goosetacob/glockchain/block"
	"github.com/goosetacob/glockchain/datastore"
	"github.com/goosetacob/glockchain/transaction"
	"github.com/sirupsen/logrus"
)

// genesisCoinbaseData
const genesisCoinbaseData = "GooseCoin `all I want is to just have fun live my life like a son of a gun`"

// Blockchain data structure
type Blockchain struct {
	leadingHash []byte
	datastore   *datastore.DataStore
}

// NewBlockchain initializes a new blockchain with a genesis block
func NewBlockchain(address string) *Blockchain {
	var err error
	datastore, err := datastore.GetDataStore()
	if err != nil {
		logrus.Error(err)
	}

	var lastHash []byte
	lastHash, err = datastore.GetLastHash()
	if err != nil {
		// TODO: come up with better way to distinguish between DB error, or new Glockchain
		coinbaseTransaction := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
		genesisBlock := newGenesisBlock(coinbaseTransaction)
		lastHash, err = datastore.AddBlock(genesisBlock)
		if err != nil {
			logrus.Error(err)
		}
	}

	return &Blockchain{lastHash, datastore}
}

// Shutdown saves and shutdsdown blockchain program
func (chain *Blockchain) Shutdown() {
	chain.datastore.Close()
}

// MineBlock mines a new block with the provided transactions
func (chain *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	lastHash, err := chain.datastore.GetLastHash()
	if err != nil {
		logrus.Panic(err)
	}

	newBlock := block.NewBlock(transactions, lastHash)

	newLeadingHash, err := chain.datastore.AddBlock(newBlock)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(newLeadingHash)
	chain.leadingHash = newLeadingHash
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

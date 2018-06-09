package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Transaction defines a transaction on the glockchain
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

const minderReward = 1000000

// NewCoinbaseTX builds a coinbase transaction ie: a special transaction which has no inputs and creates outputs ("coins")
// This is how we reward miners
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{minderReward, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.setID()

	return &tx
}

// SetID sets ID of a transaction
func (tx *Transaction) setID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		logrus.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// NewUTXOTransaction does something
func NewUTXOTransaction(from, to string, amount int, accumulated int, unspentOutputs map[string][]int) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	if accumulated < amount {
		logrus.Errorf("%v is too broke to give %v dat crypto", from, to)
		os.Exit(1)
	}

	// Build a list of inputs
	for txid, outs := range unspentOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			logrus.Errorf("error decoding Trascation ID: %v", err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if accumulated > amount {
		outputs = append(outputs, TXOutput{accumulated - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return &tx
}

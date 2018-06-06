package main

import (
	"fmt"
	"strconv"

	"github.com/goosetaco/glockchain/block"
	"github.com/goosetaco/glockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock("Send 23 BTC to Brenda")
	bc.AddBlock("Send 19 BTC to Pedro")
	bc.AddBlock("Send 11 BTC to Julio")

	for _, bc := range bc.GetBlocks() {
		fmt.Printf("Prev. hash: %x\n", bc.PrevBlockHash)
		fmt.Printf("Data: %s\n", bc.Data)
		fmt.Printf("Hash: %x\n", bc.Hash)

		pow := block.NewProofOfWork(bc)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
	}
}

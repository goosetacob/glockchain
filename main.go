package main

import (
	"flag"
	"os"

	"github.com/goosetacob/glockchain/block"
	"github.com/goosetacob/glockchain/blockchain"
	"github.com/sirupsen/logrus"
)

// CLI to operate on Blockchain
type CLI struct {
	chain *blockchain.Blockchain
}

// Run CLI
func (cli *CLI) Run() {

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addData := addCmd.String("data", "", "block data")

	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			if err := addCmd.Parse(os.Args[2:]); err != nil {
				logrus.Errorf("error: %v", err)
			} else {
				if *addData == "" {
					addCmd.Usage()
					os.Exit(1)
				}
				cli.addBlock(*addData)
			}
		case "print":
			if err := printChainCmd.Parse(os.Args[2:]); err != nil {
				logrus.Errorf("error: %v", err)
			} else {
				cli.printChain()
			}
		default:
			os.Exit(1)
		}
	} else {
		addCmd.Usage()
		printChainCmd.Usage()
		os.Exit(1)
	}
}

// addBlock adds a block to the chain
func (cli *CLI) addBlock(data string) {
	cli.chain.AddBlock(data)
	logrus.Println("Done!")
}

// printChain prints the contents of the chain
func (cli *CLI) printChain() {
	itr := cli.chain.NewIterator()

	for {
		// get the next block
		nextBlock := itr.Next()

		// print it's info
		pow := block.NewProofOfWork(nextBlock)
		logrus.Printf("\nPrevious Hash: %x\nData: %s\nHash: %x\nValidate: %t", nextBlock.PrevBlockHash, nextBlock.Data, nextBlock.Hash, pow.Validate())

		// check if we're done
		if len(nextBlock.PrevBlockHash) == 0 {
			break
		}
	}
}

func main() {
	chain := blockchain.NewBlockchain()
	defer chain.Shutdown()

	cli := CLI{chain}
	cli.Run()
}

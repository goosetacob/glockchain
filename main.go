package main

import (
	"flag"
	"os"

	"github.com/goosetacob/glockchain/blockchain"
	"github.com/goosetacob/glockchain/transaction"
	"github.com/sirupsen/logrus"
)

// CLI to operate on Blockchain
type CLI struct{}

// Run CLI
func (cli *CLI) Run() {
	//printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	balanceAddress := balanceCmd.String("address", "", "The address to send genesis block reward to")

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createAddress := createCmd.String("address", "", "The address to send genesis block reward to")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		// case "print":
		// 	if err := printChainCmd.Parse(os.Args[2:]); err != nil {
		// 		logrus.Errorf("error: %v", err)
		// 	} else {
		// 		cli.printChain()
		// 	}
		case "send":
			if err := sendCmd.Parse(os.Args[2:]); err != nil {
				logrus.Panic(err)
			} else {
				if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
					sendCmd.Usage()
					os.Exit(1)
				}
				cli.send(*sendFrom, *sendTo, *sendAmount)
			}
		case "balance":
			if err := balanceCmd.Parse(os.Args[2:]); err != nil {
				logrus.Panic(err)
			} else {
				if *balanceAddress == "" {
					balanceCmd.Usage()
					os.Exit(1)
				}
				cli.getBalance(*balanceAddress)
			}
		case "create":
			if err := createCmd.Parse(os.Args[2:]); err != nil {
				logrus.Panic(err)
			} else {
				if *createAddress == "" {
					createCmd.Usage()
					os.Exit(1)
				}
				cli.newGlockchain(*createAddress)
			}
		default:

			os.Exit(1)
		}
	} else {
		balanceCmd.Usage()
		createCmd.Usage()
		sendCmd.Usage()
		os.Exit(1)
	}
}

// printChain prints the contents of the chain
// func (cli *CLI) printChain() {
// 	itr := cli.chain.NewIterator()

// 	for {
// 		// get the next block
// 		nextBlock := itr.Next()

// 		// print it's info
// 		pow := block.NewProofOfWork(nextBlock)
// 		logrus.Printf("\nPrevious Hash: %x\nData: %s\nHash: %x\nValidate: %t", nextBlock.PrevBlockHash, nextBlock.Data, nextBlock.Hash, pow.Validate())

// 		// check if we're done
// 		if len(nextBlock.PrevBlockHash) == 0 {
// 			break
// 		}
// 	}
// }

// creates a Blockchain
func (cli *CLI) newGlockchain(address string) {
	blockchain := blockchain.NewBlockchain(address)
	defer blockchain.Shutdown()
	logrus.Println("Done!")
}

// getBalance
func (cli *CLI) getBalance(address string) {
	chain := blockchain.NewBlockchain(address)
	defer chain.Shutdown()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	logrus.Printf("Balance of '%s': %d", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	chain := blockchain.NewBlockchain(from)
	defer chain.Shutdown()

	accumulated, unspentOutputs := chain.FindSpendableOutputs(from, amount)
	tx := transaction.NewUTXOTransaction(from, to, amount, accumulated, unspentOutputs)
	chain.MineBlock([]*transaction.Transaction{tx})
	logrus.Println("Success!")
}

func main() {
	cli := CLI{}
	cli.Run()
}

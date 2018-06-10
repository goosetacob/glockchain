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
				cli.newBlockchain(*createAddress)
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

// creates a Blockchain
func (cli *CLI) newBlockchain(address string) {
	blockchain := blockchain.NewBlockchain(address)
	defer blockchain.Shutdown()
	logrus.Print("Done!")
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

// send currency from one user to another
func (cli *CLI) send(from, to string, amount int) {
	chain := blockchain.NewBlockchain(from)
	defer chain.Shutdown()

	accumulated, unspentOutputs := chain.FindSpendableOutputs(from, amount)
	tx := transaction.NewUTXOTransaction(from, to, amount, accumulated, unspentOutputs)
	chain.MineBlock([]*transaction.Transaction{tx})
	logrus.Print("Success!")
}

func main() {
	cli := CLI{}
	cli.Run()
}

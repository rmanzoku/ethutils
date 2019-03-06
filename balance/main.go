package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rmanzoku/ethutils/utils"
	"golang.org/x/net/context"
)

var (
	rpc     = os.Getenv("RPC")
	address = ""
)

func run() (err error) {
	cli, err := utils.NewEthClient(rpc)
	if err != nil {
		return
	}

	account := common.HexToAddress(address)
	balance, err := cli.BalanceAt(context.TODO(), account, nil)
	if err != nil {
		return
	}

	fmt.Println(utils.ToEther(balance).String())
	return
}

func main() {
	flag.Parse()
	address = flag.Arg(0)
	if err := run(); err != nil {
		panic(err)
	}
}

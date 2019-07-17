package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rmanzoku/ethutils/utils"
)

var (
	rpc      = os.Getenv("RPC")
	keystore = os.Getenv("ETH_KEY_PATH")
	password = os.Getenv("ETH_KEY_PASS")

	gasPriceLimit, _ = utils.GweiToWei(30)
	makeFast, _      = utils.GweiToWei(3)
)

func run(ctx context.Context, args []string) (err error) {
	cli, err := utils.NewEthClient(rpc)
	if err != nil {
		return
	}

	transactor, err := utils.NewTransactorFromKeystore(keystore, password)
	if err != nil {
		return
	}

	transactor.GasPrice, err = utils.GetGasPrice(cli, makeFast, gasPriceLimit)
	if err != nil {
		return
	}

	ret, err := utils.CancelTx(cli, transactor)
	if err != nil {
		return
	}
	fmt.Println(ret.Nonce(), ret.Hash().String())
	return
}

func main() {
	flag.Parse()
	if err := run(context.TODO(), flag.Args()); err != nil {
		panic(err)
	}
}

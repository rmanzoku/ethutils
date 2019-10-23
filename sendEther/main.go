package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rmanzoku/ethutils/utils"
	"golang.org/x/net/context"
)

var (
	rpc      = os.Getenv("RPC")
	keystore = os.Getenv("ETH_KEY_PATH")
	password = os.Getenv("ETH_KEY_PASS")
	to       = ""
	amount   = big.NewInt(0)
)

func run() (err error) {
	cli, err := utils.NewEthClient(rpc)
	if err != nil {
		return
	}

	transactor, err := utils.NewTransactorFromKeystore(keystore, password)
	if err != nil {
		return
	}
	transactor.GasPrice, _ = new(big.Int).SetString("5000000000", 10)

	fromBalance, err := cli.BalanceAt(context.TODO(), transactor.From, nil)
	if err != nil {
		return
	}

	toAddress := common.HexToAddress(to)
	toBalance, err := cli.BalanceAt(context.TODO(), toAddress, nil)
	if err != nil {
		return
	}

	if amount.Cmp(fromBalance) == 1 {
		gas, _ := new(big.Int).SetString("21000", 10)
		txFee := new(big.Int).Mul(gas, transactor.GasPrice)
		amount = new(big.Int).Sub(fromBalance, txFee)
	}

	// tx := new(types.Transaction)
	tx, err := utils.SendEther(cli, transactor, toAddress, amount)
	if err != nil {
		return
	}
	logFormat := "From:%s\tFromBalance:%v\tTo:%s\tToBalance:%v\tAmount:%v\ttx:%s\n"
	fmt.Printf(logFormat, transactor.From.String(), utils.ToEther(fromBalance), toAddress.String(), utils.ToEther(toBalance), utils.ToEther(amount).String(), tx.Hash().String())
	return
}

func main() {
	flag.Parse()
	to = flag.Arg(0)
	amountFloat, ok := new(big.Float).SetString(flag.Arg(1))
	if !ok {
		panic("Invalid arg")
	}

	var err error
	amount, err = utils.ToWei(amountFloat)
	if err != nil {
		panic(err)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

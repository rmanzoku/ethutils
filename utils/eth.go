package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func NewEthClient(rpc string) (*ethclient.Client, error) {
	conn, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	_, err = conn.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewTransactor(keystore, passphrase string) (*bind.TransactOpts, error) {
	key, err := ioutil.ReadFile(keystore)
	if err != nil {
		return nil, err
	}
	tx, err := bind.NewTransactor(bytes.NewReader(key), passphrase)
	if err != nil {
		return nil, err
	}
	tx.Context = context.TODO()
	return tx, nil
}

func SendEther(client *ethclient.Client, transactOpts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	ctx := transactOpts.Context
	nonce, err := client.NonceAt(ctx, transactOpts.From, nil)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(nonce, to, amount, 21000, transactOpts.GasPrice, nil)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	tx, err = transactOpts.Signer(types.NewEIP155Signer(chainID), transactOpts.From, tx)
	if err != nil {
		return nil, err
	}

	return tx, client.SendTransaction(transactOpts.Context, tx)
}

func ToEther(wei *big.Int) *big.Float {
	ether, _ := new(big.Float).SetString("1000000000000000000")

	w := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(w, ether)
}

func ToWei(ether *big.Float) (*big.Int, error) {
	etherDecimal, err := decimal.NewFromString(ether.String())
	if err != nil {
		return nil, err
	}

	baseWei, _ := decimal.NewFromString("1000000000000000000")

	retDecimal := etherDecimal.Mul(baseWei)
	ret, ok := new(big.Int).SetString(retDecimal.String(), 10)
	if !ok {
		return nil, errors.New("Invalit number")
	}

	return ret, nil
}

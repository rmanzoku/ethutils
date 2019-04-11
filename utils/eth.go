package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
)

var AverageBlockGenerationTime = int64(15)
var TrialCount = 10

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

func NewTransactorFromKeystore(keystore, passphrase string) (*bind.TransactOpts, error) {
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

func NewTransactorFromECDSA(filePath string) (*bind.TransactOpts, error) {
	key, err := crypto.LoadECDSA(filePath)
	if err != nil {
		return nil, err
	}
	tx := bind.NewKeyedTransactor(key)
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

func BlockByTime(client *ethclient.Client, t time.Time) (*types.Block, error) {
	ctx := context.TODO()
	result := new(types.Block)

	latest, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	latestTime := time.Unix(latest.Time().Int64(), 0)

	diffTime := latestTime.Unix() - t.Unix()
	if diffTime < 0 {
		return nil, errors.New("Block is not generated yet")
	}

	diffBlockNum := big.NewInt(diffTime / AverageBlockGenerationTime)
	targetBlockNum := new(big.Int).Sub(latest.Number(), diffBlockNum)
	cnt := 0
	for {
		if cnt > TrialCount {
			return nil, errors.New("Trial count exceeded")
		}

		eg := errgroup.Group{}
		var targetBlock *types.Block
		eg.Go(func() error {
			targetBlock, err = client.BlockByNumber(ctx, targetBlockNum)
			return err
		})
		var targetBlockBeforeOne *types.Block
		eg.Go(func() error {
			targetBlockBeforeOne, err = client.BlockByNumber(ctx, new(big.Int).Sub(targetBlockNum, big.NewInt(1)))
			return err
		})

		if err := eg.Wait(); err != nil {
			return nil, err
		}

		if targetBlockBeforeOne.Time().Int64() < t.Unix() && targetBlock.Time().Int64() >= t.Unix() {
			result = targetBlock
			break
		}

		diffTime = targetBlock.Time().Int64() - t.Unix()

		// fmt.Println(targetBlockNum, diffTime)
		diffBlockNum = big.NewInt(diffTime/AverageBlockGenerationTime/2 + 1)
		targetBlockNum = new(big.Int).Sub(targetBlock.Number(), diffBlockNum)
		cnt++
	}

	return result, nil
}

func abs(a int64) int64 {
	if a < 0 {
		return -1 * a
	}
	return a
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

package main

import (
	"flag"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rmanzoku/ethutils/ecrecover"
)

func run(args []string) (err error) {
	key, err := crypto.LoadECDSA("key")
	if err != nil {
		return err
	}
	transactor := bind.NewKeyedTransactor(key)

	message := args[0]
	msg, _ := ecrecover.ToEthSignedMessageHash([]byte(message))
	sig, _ := crypto.Sign(msg, key)

	fmt.Println("message:", message)
	fmt.Println("signer:", transactor.From.String())
	fmt.Println("sig:", hexutil.Encode(sig))
	return nil
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		panic(err)
	}
}

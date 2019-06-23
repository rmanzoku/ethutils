package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

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
	if !strings.HasPrefix(message, "0x") {
		message = "0x" + hex.EncodeToString([]byte(message))
	}

	h, err := hexutil.Decode(message)
	if err != nil {
		return
	}
	msg := ecrecover.ToEthSignedMessageHash(h)
	if err != nil {
		return
	}

	sig, err := crypto.Sign(msg, key)
	if err != nil {
		return
	}
	sig[64] += 27

	fmt.Println(sig)
	signer, err := ecrecover.Recover(msg, sig)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("arg:", args[0])
	fmt.Println("message:", message)
	// fmt.Println("message sha3 hex:", hexutil.Encode(keccakMsg))
	fmt.Println("message hash hex:", hexutil.Encode(msg))
	fmt.Println("signer:", transactor.From.String())
	fmt.Println("recover:", signer.String())
	fmt.Println("sig:", hexutil.Encode(sig))
	return nil
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		panic(err)
	}
}

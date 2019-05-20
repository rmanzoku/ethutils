package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func run(args []string) (err error) {
	file := args[0]

	key, err := crypto.LoadECDSA(file)
	if err != nil {
		return
	}

	pubkey := key.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pubkey)
	fmt.Print(addr.String())
	return
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		panic(err)
	}
}

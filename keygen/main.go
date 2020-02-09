package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func run() (err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	fmt.Println(string(crypto.FromECDSA(key)))
	k := hex.EncodeToString(crypto.FromECDSA(key))
	fmt.Print(k)
	return
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

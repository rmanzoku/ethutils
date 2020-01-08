package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

func run(args []string) (err error) {

	j, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scanner.Text()

	k, err := keystore.DecryptKey(j, scanner.Text())
	if err != nil {
		return
	}

	b := crypto.FromECDSA(k.PrivateKey)
	fmt.Println(hex.EncodeToString(b))
	return
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		panic(err)
	}
}

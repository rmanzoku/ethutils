package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	flag.Parse()
	times := 0
	if len(flag.Args()) == 0 {
		times = 1
	} else {
		times, _ = strconv.Atoi(flag.Args()[0])
	}

	ret := make([]string, times)
	for i := 0; i < times; i++ {
		seed := time.Now().UnixNano()
		r := rand.New(rand.NewSource(seed))
		b := big.NewInt(r.Int63())
		k := crypto.Keccak256(b.Bytes())
		h := hex.EncodeToString(k)
		ret[i] = "\"0x" + h + "\""
	}

	fmt.Println("[" + strings.Join(ret, ", ") + "]")
}

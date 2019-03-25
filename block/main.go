package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rmanzoku/ethutils/utils"
)

var (
	rpc = os.Getenv("RPC")
	t   = time.Unix(1553307665, 0)
)

func run() (err error) {
	cli, err := utils.NewEthClient(rpc)
	if err != nil {
		return
	}

	block, err := utils.BlockByTime(cli, t)
	if err != nil {
		return
	}

	fmt.Println(block.Time(), t.Unix())
	return
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}

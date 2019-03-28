package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rmanzoku/ethutils/utils"
)

var (
	rpc = os.Getenv("RPC")
	t   = time.Now()
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
	tInt, err := strconv.ParseInt(flag.Arg(0), 10, 64)
	if err != nil {
		panic(err)
	}
	t = time.Unix(tInt, 0)
	if err := run(); err != nil {
		panic(err)
	}
}

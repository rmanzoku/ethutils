package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/utils"
	"golang.org/x/net/context"
)

var (
	rpc = os.Getenv("RPC")
)

func run(args []string) (err error) {
	ctx := context.TODO()

	target, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return
	}
	cli, err := utils.NewEthClient(rpc)
	if err != nil {
		return
	}

	latestBlock, err := cli.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}

	latest := latestBlock.Number().Int64()
	now := time.Now().Unix()

	diff := target - latest
	if diff < 0 {
		return errors.New("past block target")
	}

	fmt.Printf("target: %d, current: %d, left: %d\n", target, latest, diff)

	time15 := time.Unix((diff*15)+now, 0)
	time13 := time.Unix((diff*13)+now, 0)
	time10 := time.Unix((diff*10)+now, 0)

	fmt.Printf("time15: " + time15.Format(time.RFC1123Z) + "\n")
	fmt.Printf("time13: " + time13.Format(time.RFC1123Z) + "\n")
	fmt.Printf("time10: " + time10.Format(time.RFC1123Z) + "\n")

	return
}

func main() {
	flag.Parse()
	if err := run(flag.Args()); err != nil {
		panic(err)
	}
}

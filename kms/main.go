package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/utils"
	"golang.org/x/net/context"
)

var (
	rpc     = os.Getenv("RPC")
	profile = os.Getenv("AWS_PROFILE")
	region  = os.Getenv("AWS_DEFAULT_REGION")
	to      = "0xd868711BD9a2C6F1548F5f4737f71DA67d821090"
	keyID   = "465aaa3c-e03b-41e1-b675-6206c1227f83"
	amount  = big.NewInt(0)
)

func NewKMSTransactor(k *kms.KMS, id string) (*bind.TransactOpts, error) {
	in := &kms.GetPublicKeyInput{
		KeyId: aws.String(id),
	}
	out, err := k.GetPublicKey(in)
	if err != nil {
		return nil, err
	}

	pubkey, err := crypto.UnmarshalPubkey(out.PublicKey)
	if err != nil {
		return
	}
	keyAddr := crypto.PubkeyToAddress(*pubkey)

	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}

			in := &kms.SignInput{
				KeyId:            aws.String(id),
				Message:          signer.Hash(tx).Bytes(),
				SigningAlgorithm: aws.String("ECDSA_SHA_256"),
				MessageType:      aws.String("DIGEST"),
			}
			out, err := k.Sign(in)
			if err != nil {
				return nil, err
			}

			signature := out.Signature
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

func run() (err error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	})
	if err != nil {
		return
	}
	k := kms.New(sess)

	transactor, err := NewKMSTransactor(k, keyID)
	if err != nil {
		return
	}
	transactor.Context = context.TODO()
	cli, err := utils.NewEthClient(rpc)
	if err != nil {
		return
	}

	toAddress := common.HexToAddress(to)
	tx, err := utils.SendEther(cli, transactor, toAddress, amount)
	if err != nil {
		return
	}
	fmt.Println(tx.Hash().String())
	return

}

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

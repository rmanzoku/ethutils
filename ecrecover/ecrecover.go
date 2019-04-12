package ecrecover

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func Example() {
	var signature string
	var message []byte

	sig, err := hexutil.Decode(signature)
	if err != nil {
		return
	}

	hash, err := ToEthSignedMessageHash(message)
	if err != nil {
		return
	}

	signer, err := Recover(hash, sig)
	if err != nil {
		return
	}
	fmt.Println(signer)
}

func Recover(message []byte, sig []byte) (common.Address, error) {
	nilAddress := common.HexToAddress("0x0")

	if sig[64] != 27 && sig[64] != 28 {
		return nilAddress, errors.New("recovery error")
	}
	sig[64] -= 27

	p, err := crypto.SigToPub(message, sig)
	if err != nil {
		return nilAddress, err
	}

	return crypto.PubkeyToAddress(*p), nil
}

func ToEthSignedMessageHash(hash []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, errors.New("hash len is not 32")
	}
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	return crypto.Keccak256([]byte(msg)), nil
}

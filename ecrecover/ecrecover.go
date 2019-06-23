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

	hash := ToEthSignedMessageHash(message)

	signer, err := Recover(hash, sig)
	if err != nil {
		return
	}
	fmt.Println(signer)
}

func Recover(message []byte, sig []byte) (common.Address, error) {
	nilAddress := common.HexToAddress("0x0")

	if len(sig) < 63 {
		return nilAddress, errors.New("sig size is not under 63")
	}
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

func ToEthSignedMessageHash(message []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	fmt.Println(msg)
	return Keccak256([]byte(msg))
}

func Keccak256(data []byte) []byte {
	return crypto.Keccak256(data)
}

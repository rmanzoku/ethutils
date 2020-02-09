package ethersign

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Signature []byte

func IsHexSignature(s string) bool {
	if s[0:2] != "0x" {
		return false
	}
	length := 132
	if len(s) != length {
		return false
	}
	if !isHex(s[2:]) {
		return false
	}
	return true
}

func ParseSignature(signature string) (Signature, error) {
	if !IsHexSignature(signature) {
		return nil, errors.New("Invalid signature: " + signature)
	}
	sig, err := decodeHex(signature)
	if err != nil {
		return nil, err
	}
	return Signature(sig), nil
}

func ParseSignatureMust(signature string) Signature {
	ret, err := ParseSignature(signature)
	if err != nil {
		panic(err)
	}
	return ret
}

func (s Signature) Bytes() []byte {
	return []byte(s)
}

func (s Signature) String() string {
	return encodeToHex(s)
}

func ToEthSignedMessageHash(message []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return Keccak256([]byte(msg))
}

func Keccak256(data []byte) []byte {
	return crypto.Keccak256(data)
}

func encodeToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func decodeHex(s string) ([]byte, error) {
	if s[0:2] != "0x" {
		return nil, errors.New("hex must start with 0x")
	}
	return hex.DecodeString(s[2:])
}

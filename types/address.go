package types

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Address string

// BlackHoleAddress is 0x0000000000000000000000000000000000000000
var BlackHoleAddress = Address("0x0000000000000000000000000000000000000000")

func IsHexAddress(s string) bool {
	if s == "0x0" {
		return true
	}
	if s[0:2] != "0x" {
		return false
	}
	addressLength := 2 + 40
	if len(s) != addressLength {
		return false
	}
	if !isHex(s[2:]) {
		return false
	}
	return true
}

func ParseAddress(address string) (Address, error) {
	if !IsHexAddress(address) {
		return "", errors.New("Invalid address: " + address)
	}
	if address == "0x0" {
		return BlackHoleAddress, nil
	}
	return Address(strings.ToLower(address)), nil
}

func ParseAddressMust(address string) Address {
	ret, err := ParseAddress(address)
	if err != nil {
		panic(err)
	}
	return ret
}

func (a Address) String() string {
	return string(a)
}

func (a Address) Address() common.Address {
	return common.HexToAddress(a.String())
}

func (a Address) IsBlackHole() bool {
	if a.String() == BlackHoleAddress.String() {
		return true
	}
	return false
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

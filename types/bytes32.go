package types

import (
	"encoding/hex"
	"errors"
	"strings"
)

type Bytes32 string

// ZeroBytes32 is 0x0000000000000000000000000000000000000000000000000000000000000000
var ZeroBytes32 = Bytes32("0x0000000000000000000000000000000000000000000000000000000000000000")

func IsHexByte32(s string) bool {
	if s == "0x0" {
		return true
	}
	if s[0:2] != "0x" {
		return false
	}
	length := 2 + 64
	if len(s) != length {
		return false
	}
	if !isHex(s[2:]) {
		return false
	}
	return true
}

func ParseBytes32(bytes32 string) (Bytes32, error) {
	if !IsHexByte32(bytes32) {
		return "", errors.New("Invalid bytes32: " + bytes32)
	}
	if bytes32 == "0x0" {
		return ZeroBytes32, nil
	}
	return Bytes32(strings.ToLower(bytes32)), nil
}

func ParseBytes32Must(bytes32 string) Bytes32 {
	ret, err := ParseBytes32(bytes32)
	if err != nil {
		panic(err)
	}
	return ret
}

func (b Bytes32) String() string {
	return string(b)
}

func (b Bytes32) Bytes32() [32]byte {
	bytes, err := hex.DecodeString(b.String()[2:])
	if err != nil {
		panic(err)
	}
	var ret [32]byte
	copy(ret[:], bytes[0:32])
	return ret
}

func (b Bytes32) IsZero() bool {
	if b.String() == ZeroBytes32.String() {
		return true
	}
	return false
}

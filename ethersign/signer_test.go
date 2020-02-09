package ethersign

import (
	"os"
	"testing"
)

func TestSigner(t *testing.T) {
	msg := []byte("Hello world")
	is := initTesting(t)
	signer, err := NewSignerFromHex(os.Getenv("PRIVATE_KEY"))
	is.Nil(err)
	print(signer.Address())

	sig, err := signer.EthereumSign(msg)
	is.Nil(err)
	print(sig.String())
	hash := ToEthSignedMessageHash(msg)

	addr, err := RecoveryByCrypto(hash, sig)
	is.Nil(err)
	print(addr)
}

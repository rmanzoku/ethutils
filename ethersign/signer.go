package ethersign

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	privKey *ecdsa.PrivateKey
}

func NewSigner(key *ecdsa.PrivateKey) (*Signer, error) {
	return &Signer{
		privKey: key,
	}, nil
}

func NewSignerFromHex(hex string) (*Signer, error) {
	b, err := decodeHex(hex)
	if err != nil {
		return nil, err
	}
	k, err := crypto.ToECDSA(b)
	if err != nil {
		return nil, err
	}
	return NewSigner(k)
}

func (s *Signer) Public() ecdsa.PublicKey {
	return *s.privKey.Public().(*ecdsa.PublicKey)
}

func (s *Signer) Address() Address {
	return ParseAddressMust(crypto.PubkeyToAddress(s.Public()).String())
}

func (s *Signer) EthereumSign(msg []byte) (Signature, error) {
	h := ToEthSignedMessageHash(msg)

	sig, err := crypto.Sign(h, s.privKey)
	if err != nil {
		return nil, err
	}
	if sig[64] < 27 {
		sig[64] += 27
	}

	return ParseSignature(encodeToHex(sig))
}

func RecoveryByCrypto(hash []byte, sig Signature) (Address, error) {
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	pub, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return BlackHoleAddress, err
	}
	addr := crypto.PubkeyToAddress(*pub)
	return ParseAddress(addr.String())
}

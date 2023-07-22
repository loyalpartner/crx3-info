package crx3

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/sha256"
	"hash"
)

type ECDSAVerifier struct {
	algorithm crypto.Hash
	publicKey *ecdsa.PublicKey
	signature []byte
	hasher    hash.Hash
}

func NewECDSAVerifier(
	algorithm crypto.Hash,
	publicKey *ecdsa.PublicKey,
	signature []byte,
) SignatureVerifier {
	return &ECDSAVerifier{
		algorithm: algorithm,
		publicKey: publicKey,
		signature: signature,
	}
}

func (v *ECDSAVerifier) Hasher() hash.Hash {
	if v.hasher == nil {
		v.hasher = sha256.New()
	}
	return v.hasher
}

func (v *ECDSAVerifier) Update(data []byte) error {
	_, err := v.Hasher().Write(data)
	return err
}

func (v *ECDSAVerifier) Verify() error {
	hash := v.Hasher().Sum(nil)
	if !ecdsa.VerifyASN1(v.publicKey, hash, v.signature) {
		return ErrInvalidSignature
	}
	return nil
}

package crx3

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"hash"

	"github.com/loyalpartner/crx3-info/pb"
)

type ECDSAVerifier struct {
	algorithm crypto.Hash
	publicKey *ecdsa.PublicKey
	signature []byte
	hasher    hash.Hash
}

func NewECDSAVerifier(
	algorithm crypto.Hash,
	proof *pb.AsymmetricKeyProof,
) (SignatureVerifier, error) {

	pubilcKey, err := x509.ParsePKIXPublicKey(proof.PublicKey)
	if err != nil {
		return nil, err
	}

	ecdsaPublicKey, ok := pubilcKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrNotEcdsaPublicKey
	}
	return &ECDSAVerifier{
		algorithm: algorithm,
		publicKey: ecdsaPublicKey,
		signature: proof.Signature,
	}, nil
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

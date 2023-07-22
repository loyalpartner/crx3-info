package crx3

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"

	"github.com/loyalpartner/crx3-info/pb"
)

type SignatureAlgorithm int

const (
	RSA_PKCS1_SHA1 SignatureAlgorithm = iota
	RSA_PKCS1_SHA256
	ECDSA_SHA256
	// RSA_PSS_SHA256
)

type SignatureVerifier interface {
	Update(data []byte) error
	Verify() error
}

func NewSignatureVerifier(algorithm crypto.Hash,
	proof *pb.AsymmetricKeyProof) (SignatureVerifier, error) {

	pubKey, err := x509.ParsePKIXPublicKey(proof.PublicKey)
	if err != nil {
		return nil, err
	}

	switch pubKey.(type) {
	case *rsa.PublicKey:
		rsaPubKey, _ := pubKey.(*rsa.PublicKey)
		return NewRSAVerifier(crypto.SHA256, rsaPubKey, proof.Signature), nil
	case *ecdsa.PublicKey:
		ecdsaPubKey, _ := pubKey.(*ecdsa.PublicKey)
		return NewECDSAVerifier(crypto.SHA256, ecdsaPubKey, proof.Signature), nil
	default:
		return nil, ErrUnsupportedAlgorithm
	}
}

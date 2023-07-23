package crx3

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

type RSAVerifier struct {
	algorithm crypto.Hash
	publicKey *rsa.PublicKey
	signature []byte
	hasher    hash.Hash
}

func NewRSAVerifier(algorithm crypto.Hash, pubKey *rsa.PublicKey, signature []byte) Verifier {

	return &RSAVerifier{
		publicKey: pubKey,
		algorithm: algorithm,
		signature: signature,
	}
}

func (s *RSAVerifier) Hasher() hash.Hash {
	if s.hasher == nil {
		switch s.algorithm {
		case crypto.SHA256:
			s.hasher = sha256.New()
		default:
			s.hasher = sha1.New()
		}
	}
	return s.hasher
}

func (s *RSAVerifier) Update(data []byte) error {
	_, err := s.Hasher().Write(data)
	return err
}

func (s *RSAVerifier) Verify() error {
	hash := s.Hasher().Sum(nil)
	return rsa.VerifyPKCS1v15(s.publicKey,
		s.algorithm,
		hash,
		s.signature)
}

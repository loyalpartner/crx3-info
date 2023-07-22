package crx3

import "errors"

var (
	ErrInvalidSignature     = errors.New("invalid signature")
	ErrNotRSAPublicKey      = errors.New("not rsa public key")
	ErrNotEcdsaPublicKey    = errors.New("not ecdsa public key")
	ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")
)

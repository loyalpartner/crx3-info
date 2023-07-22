package crx3

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

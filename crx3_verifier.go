package crx3

import (
	"crypto"
)

const (
	SignedDataPrefix string = "CRX3 SignedData\x00"
)

type verifier struct {
	crx       *Crx3
	verifiers []SignatureVerifier
}

type verifierList []SignatureVerifier

func (v *verifierList) Append(verifier SignatureVerifier) {
	*v = append(*v, verifier)
}

func NewVerifier(crx *Crx3) (*verifier, error) {
	v := &verifier{}
	if err := v.initialize(crx); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *verifier) initialize(crx *Crx3) error {
	c.crx = crx
	for _, rsaProof := range crx.Header.Sha256WithRsa {
		rsaVerifier, err := NewRSAVerifier(crypto.SHA256, rsaProof)
		if err != nil {
			return err
		}
		c.verifiers = append(c.verifiers, rsaVerifier)
	}
	for _, ecdsaProof := range crx.Header.Sha256WithEcdsa {
		edcsaVerifier, err := NewECDSAVerifier(crypto.SHA256, ecdsaProof)
		if err != nil {
			return nil
		}
		c.verifiers = append(c.verifiers, edcsaVerifier)
	}
	return nil
}

func (c *verifier) update(data []byte) {
	for _, v := range c.verifiers {
		v.Update(data)
	}
}

func (c *verifier) verify() error {
	for _, verifier := range c.verifiers {
		if err := verifier.Verify(); err != nil {
			return err
		}
	}
	return nil
}

func (c *verifier) Verify() error {
	hdr := c.crx.Header

	c.update([]byte(SignedDataPrefix))
	c.update(c.crx.LeEncodedSignedDataLen())
	c.update(hdr.SignedHeaderData)
	c.update(c.crx.Archive)

	return c.verify()
}

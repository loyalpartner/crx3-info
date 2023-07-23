package crx3

import (
	"crypto"

	"github.com/loyalpartner/crx3-info/pb"
	"golang.org/x/sync/errgroup"
)

const (
	SignedDataPrefix string = "CRX3 SignedData\x00"
)

type Verifier interface {
	Update(data []byte) error
	Verify() error
}

type Crx3Verifier struct {
	crx       *Crx3
	verifiers []Verifier
}

func (c *Crx3Verifier) Verify() error {
	crx := c.crx

	c.Update([]byte(SignedDataPrefix))
	c.Update(crx.LeEncodedSignedDataLen())
	c.Update(crx.Header.SignedHeaderData)
	c.Update(crx.Archive)

	return c.verify()
}

func (c *Crx3Verifier) Update(data []byte) {
	for _, v := range c.verifiers {
		v.Update(data)
	}
}

func (c *Crx3Verifier) addVerifiers(proofs []*pb.AsymmetricKeyProof) error {
	g := &errgroup.Group{}
	for _, proof := range proofs {
		p := proof
		g.Go(func() error {
			verifier, err := NewVerifierImpl(crypto.SHA256, p)
			if err == nil {
				c.verifiers = append(c.verifiers, verifier)
			}
			return err
		})
	}
	return g.Wait()
}

func (c *Crx3Verifier) verify() error {
	g := &errgroup.Group{}
	for _, verifier := range c.verifiers {
		g.Go(verifier.Verify)
	}
	return g.Wait()
}

func Verify(crx *Crx3) error {
	verifier, err := NewVerifier(crx)
	if err != nil {
		return err
	}
	return verifier.Verify()
}

func NewVerifier(crx *Crx3) (*Crx3Verifier, error) {
	v := &Crx3Verifier{crx: crx}

	g := &errgroup.Group{}
	g.Go(func() error { return v.addVerifiers(crx.Header.Sha256WithRsa) })
	g.Go(func() error { return v.addVerifiers(crx.Header.Sha256WithEcdsa) })

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return v, nil
}

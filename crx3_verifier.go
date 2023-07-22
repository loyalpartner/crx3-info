package crx3

import (
	"crypto"

	"github.com/loyalpartner/crx3-info/pb"
	"golang.org/x/sync/errgroup"
)

const (
	SignedDataPrefix string = "CRX3 SignedData\x00"
)

type TaskFunc func() error

type crxVerifier struct {
	crx       *Crx3
	verifiers []SignatureVerifier
}

func NewVerifier(crx *Crx3) (*crxVerifier, error) {
	v := &crxVerifier{}
	if err := v.initialize(crx); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *crxVerifier) initialize(crx *Crx3) error {
	c.crx = crx

	g := &errgroup.Group{}
	g.Go(c.fnInitializeVerifiers(c.crx.Header.Sha256WithRsa))
	g.Go(c.fnInitializeVerifiers(c.crx.Header.Sha256WithEcdsa))
	return g.Wait()
}

func (c *crxVerifier) fnInitializeVerifiers(proofs []*pb.AsymmetricKeyProof) TaskFunc {
	return func() error {
		return c.initializeVerifiers(proofs)
	}
}

func (c *crxVerifier) initializeVerifiers(proofs []*pb.AsymmetricKeyProof) error {
	g := &errgroup.Group{}
	for _, proof := range proofs {
		g.Go(c.fnAddSignatureVerifier(proof))
	}
	return g.Wait()
}

func (c *crxVerifier) fnAddSignatureVerifier(proof *pb.AsymmetricKeyProof) TaskFunc {
	return func() error {
		verifier, err := NewSignatureVerifier(crypto.SHA256, proof)
		if err == nil {
			c.verifiers = append(c.verifiers, verifier)
		}
		return err
	}
}

func (c *crxVerifier) update(data []byte) {
	for _, v := range c.verifiers {
		v.Update(data)
	}
}

func (c *crxVerifier) verify() error {
	g := &errgroup.Group{}
	for _, verifier := range c.verifiers {
		g.Go(verifier.Verify)
	}
	return g.Wait()
}

func (c *crxVerifier) Verify() error {
	hdr := c.crx.Header

	// TODO: add crx_reader.go
	c.update([]byte(SignedDataPrefix))
	c.update(c.crx.LeEncodedSignedDataLen())
	c.update(hdr.SignedHeaderData)
	c.update(c.crx.Archive)

	return c.verify()
}

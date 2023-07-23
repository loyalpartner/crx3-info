package crx3

import "github.com/loyalpartner/crx3-info/pb"

type Header struct {
	Sha256WithRsa   []*pb.AsymmetricKeyProof
	Sha256WithEcdsa []*pb.AsymmetricKeyProof
	// SignedHeaderData *pb.SignedData
	SignedHeaderData []byte
}

func NewHeader(header *pb.CrxFileHeader) *Header {
	return &Header{
		Sha256WithRsa:    header.Sha256WithRsa,
		Sha256WithEcdsa:  header.Sha256WithEcdsa,
		SignedHeaderData: header.SignedHeaderData,
	}
}

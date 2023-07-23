package crx3

import (
	"encoding/binary"

	"github.com/loyalpartner/crx3-info/pb"
	"google.golang.org/protobuf/encoding/protojson"
)

// crx3 format:
// https://source.chromium.org/chromium/chromium/src/+/main:components/crx_file/crx3.proto
type Crx3 struct {
	Magic      string
	Version    int
	HeaderSize int
	Header     *Header
	PbHeader   *pb.CrxFileHeader
	Archive    []byte
	ID         string
}

func (c *Crx3) JsonEncodedHeader() string {
	return protojson.Format(c.PbHeader)
}

func (c *Crx3) LeEncodedSignedDataLen() []byte {
	l := len(c.Header.SignedHeaderData)
	return binary.LittleEndian.AppendUint32(nil, uint32(l))
}

func NewCrx3(magic string, version, headerSize int, archive []byte) *Crx3 {
	return &Crx3{
		Magic:      magic,
		Version:    version,
		HeaderSize: headerSize,
		Header:     &Header{},
		PbHeader:   &pb.CrxFileHeader{},
		Archive:    archive,
	}
}

package crx3

import (
	"encoding/binary"
	"os"

	"github.com/loyalpartner/crx3-info/pb"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// crx3 format:
// https://source.chromium.org/chromium/chromium/src/+/main:components/crx_file/crx3.proto
type Crx3 struct {
	Path       string
	Magic      string
	Version    int
	HeaderSize int
	Header     *pb.CrxFileHeader
	Archive    []byte
}

func NewCrx3(path string) *Crx3 {
	return &Crx3{Path: path}
}

func (c *Crx3) Load() error {
	raw, err := os.ReadFile(c.Path)
	if err != nil {
		return err
	}

	c.Header = &pb.CrxFileHeader{}
	c.Magic = string(raw[0:4])
	c.Version = int(binary.LittleEndian.Uint32(raw[4:8]))
	c.HeaderSize = int(binary.LittleEndian.Uint32(raw[8:12]))

	headerOffset := 12
	archiveOffset := headerOffset + c.HeaderSize

	header := raw[headerOffset:archiveOffset]
	if err = proto.Unmarshal(header, c.Header); err != nil {
		return err
	}

	c.Archive = raw[archiveOffset:]
	return nil
}

func (c *Crx3) JsonEncodedHeader() string {
	return protojson.Format(c.Header)
}

func (c *Crx3) SignedData() *pb.SignedData {
	signedData := &pb.SignedData{}
	proto.Unmarshal(c.Header.SignedHeaderData, signedData)
	return signedData
}

func (c *Crx3) LeEncodedSignedDataLen() []byte {
	l := len(c.Header.SignedHeaderData)
	return binary.LittleEndian.AppendUint32(nil, uint32(l))
}

func (c *Crx3) CrxId() string {
	if c.Header == nil {
		return ""
	}
	return FromBytes(c.SignedData().CrxId)
}

func (c *Crx3) Verify() error {
	verifier, err := NewVerifier(c)
	if err != nil {
		return err
	}
	return verifier.Verify()
}

package crx3

import (
	"encoding/binary"
	"os"

	"github.com/loyalpartner/crx3-info/pb"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	HEADER_OFFSET int = 12
)

// crx3 format:
// https://source.chromium.org/chromium/chromium/src/+/main:components/crx_file/crx3.proto
type Crx3 struct {
	Path       string
	Magic      string
	Version    int
	HeaderSize int
	Header     *pb.CrxFileHeader
	Data       []byte
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

	header := raw[HEADER_OFFSET : HEADER_OFFSET+c.HeaderSize]
	if err = proto.Unmarshal(header, c.Header); err != nil {
		return err
	}

	return nil
}

func (c *Crx3) HeaderDetails() string {
	return protojson.Format(c.Header)
}

func (c *Crx3) CrxId() string {
	if c.Header == nil {
		return ""
	}

	signedData := &pb.SignedData{}
	proto.Unmarshal(c.Header.SignedHeaderData, signedData)

	return FromBytes(signedData.CrxId)
}

package crx3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"os"

	"github.com/golang/protobuf/proto"
)

func Read(path string) (*Crx3, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	magic := string(raw[0:4])
	version := int(binary.LittleEndian.Uint32(raw[4:8]))
	headerSize := int(binary.LittleEndian.Uint32(raw[8:12]))
	headerOffset := 12
	archiveOffset := headerOffset + headerSize
	archive := raw[archiveOffset:]

	crx := NewCrx3(magic, version, headerSize, archive)

	header := raw[headerOffset:archiveOffset]
	if err = proto.Unmarshal(header, crx.PbHeader); err != nil {
		return nil, err
	}
	crx.Header = NewHeader(crx.PbHeader)

	crxId := ReadID(crx.Header.SignedHeaderData)
	crx.ID = crxId

	return crx, nil
}

func ReadID(raw []byte) string {
	hexStr := hex.EncodeToString(raw)
	return convertHexadecimalToIDAlphabet(hexStr)
}

// Converts a normal hexadecimal string into the alphabet used by extensions.
// We use the characters 'a'-'p' instead of '0'-'f' to avoid ever having a
// completely numeric host, since some software interprets that as an IP
// address.
func convertHexadecimalToIDAlphabet(hexStr string) string {
	buf := bytes.NewBuffer(nil)
	offset := rune(0)
	for _, c := range hexStr {
		if c >= 48 && c <= 57 {
			offset = 'a' - '0'
		} else if c >= 97 && c <= 102 {
			offset = 'k' - 'a'
		}
		buf.WriteRune(c + offset)
	}
	return buf.String()
}

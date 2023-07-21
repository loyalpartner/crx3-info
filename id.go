package crx3

import (
	"bytes"
	"encoding/hex"
)

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

func FromBytes(data []byte) string {
	hexStr := hex.EncodeToString(data)
	return convertHexadecimalToIDAlphabet(hexStr)
}

package crx3

import (
	"bytes"
	"fmt"
)

const SYMBOLS = "abcdefghijklmnopqrstuvwxyz"

func strIDx() map[rune]int {
	index := make(map[rune]int)
	src := "0123456789abcdef"
	for i, char := range src {
		index[char] = i
	}
	return index
}

func FromBytes(data []byte) string {
	sid := fmt.Sprintf("%x", data)
	idx := strIDx()
	buf := bytes.NewBuffer(nil)
	for _, char := range sid {
		index := idx[char]
		buf.WriteString(string(SYMBOLS[index]))
	}
	return buf.String()
}

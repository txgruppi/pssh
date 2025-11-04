package main

import (
	"bufio"
	"bytes"
)

func prefixLines(prefix []byte, lines []byte) []byte {
	out := make([]byte, 0, len(lines)*2)
	scanner := bufio.NewScanner(bytes.NewReader(lines))
	for scanner.Scan() {
		out = append(out, prefix...)
		out = append(out, scanner.Bytes()...)
		out = append(out, '\n')
	}
	return out
}

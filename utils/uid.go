package utils

import "github.com/rs/xid"

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NewID returns an unique 20 character long string
func NewID() string {
	return reverse(xid.New().String())
}

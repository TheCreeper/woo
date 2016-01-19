package urigen

import (
	"crypto/rand"
	"math/big"
)

const Symbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Generate will attempt to generate a random set of bytes of the
// specified size. There is no limit on the length.
func Generate(size int) (b []byte, err error) {
	for i := 0; i < size; i++ {
		c, err := RandChar()
		if err != nil {
			return nil, err
		}
		b = append(b, c)
	}
	return
}

// RandChar returns a random character from Symbols.
func RandChar() (c byte, err error) {
	max := big.NewInt(int64(len(Symbols)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return
	}
	c = Symbols[n.Int64()]
	return
}

package random

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var ErrRandomGenerator = errors.New("random generator error")

type Generator interface {
	Uint(min, max uint) (uint, error)
}

type Crypto struct {
}

func (c *Crypto) Uint(min, max uint) (uint, error) {
	bg := big.NewInt(int64(max - min + 1))

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, fmt.Errorf("%s | %w", err.Error(), ErrRandomGenerator)
	}

	return uint(n.Uint64()) + min, nil
}

func NewCrypto() *Crypto {
	return &Crypto{}
}

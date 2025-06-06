package random

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var ErrRandomGenerator = errors.New("random generator error")

type Generator interface {
	Uint(minValue, maxValue uint) (uint, error)
}

type Crypto struct {
}

func NewCrypto() *Crypto {
	return &Crypto{}
}

func (c *Crypto) Uint(minValue, maxValue uint) (uint, error) {
	val := maxValue - minValue + 1
	//nolint:gosec
	bg := big.NewInt(int64(val))

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, fmt.Errorf("%s | %w", err.Error(), ErrRandomGenerator)
	}

	return uint(n.Uint64()) + minValue, nil
}

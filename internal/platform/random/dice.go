package random

import (
	"errors"
	"fmt"
)

var ErrDice = errors.New("dice error")

type Dice struct {
	generator Generator
	Sides     uint
}

func NewCubeDice() *Dice {
	return &Dice{generator: NewCrypto(), Sides: uint(6)}
}

func NewDice(generator Generator, sides uint) *Dice {
	return &Dice{generator: generator, Sides: sides}
}

func (d Dice) Roll() (uint, error) {
	if d.Sides == 0 {
		return 0, fmt.Errorf("%w | must have at least one side", ErrDice)
	}

	return d.generator.Uint(1, d.Sides)
}

package random

import (
	"errors"
	"fmt"
)

var ErrDice = errors.New("dice error")
var ErrDiceNotEnoughSides = fmt.Errorf("%w: dice not enough sides", ErrDice)

const DiceSides = uint(6)

type Dice struct {
	generator Generator
	Sides     uint
}

func NewCubeDice() *Dice {
	return &Dice{generator: NewCrypto(), Sides: DiceSides}
}

func NewDice(generator Generator, sides uint) *Dice {
	return &Dice{generator: generator, Sides: sides}
}

func (d Dice) Roll() (uint, error) {
	if d.Sides == 0 {
		return 0, ErrDiceNotEnoughSides
	}

	res, err := d.generator.Uint(1, d.Sides)
	if err != nil {
		return 0, fmt.Errorf("%w: generate number: %w", ErrDice, err)
	}

	return res, nil
}

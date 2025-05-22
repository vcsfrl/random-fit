package random_test

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/platform/random"
	"testing"
)

func TestDice(t *testing.T) {
	suite.Run(t, new(DiceFixture))
}

type DiceFixture struct {
	suite.Suite
	dice      *random.Dice
	generator *MockGenerator
}

func (df *DiceFixture) SetupTest() {
	df.generator = &MockGenerator{}
	df.dice = random.NewDice(df.generator, 0)
}

func (df *DiceFixture) TestDice() {
	for i := uint(2); i < uint(100); i += 2 {
		df.dice.Sides = i

		pick, err := df.dice.Roll()
		df.Nil(err)
		df.Equal(df.generator.lastMin, uint(1))
		df.Equal(df.generator.lastMax, df.dice.Sides)
		df.Equal(pick, df.dice.Sides-1)
	}
}

func (df *DiceFixture) TestDiceErr() {
	df.dice.Sides = 0
	pick, err := df.dice.Roll()

	df.Equal(pick, uint(0))
	df.True(errors.Is(err, random.ErrDice))
	df.Equal(err.Error(), "dice error | must have at least one side")
}

type MockGenerator struct {
	lastMin, lastMax uint
}

func (m *MockGenerator) Uint(min, max uint) (uint, error) {
	m.lastMin = min
	m.lastMax = max

	return max - 1, nil
}

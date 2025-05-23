package random_test

import (
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
		df.Require().NoError(err)
		df.Equal(uint(1), df.generator.lastMin)
		df.Equal(df.generator.lastMax, df.dice.Sides)
		df.Equal(pick, df.dice.Sides-1)
	}
}

func (df *DiceFixture) TestDiceErr() {
	df.dice.Sides = 0
	pick, err := df.dice.Roll()

	df.Equal(uint(0), pick)
	df.ErrorIs(err, random.ErrDiceNotEnoughSides)
}

type MockGenerator struct {
	lastMin, lastMax uint
}

func (m *MockGenerator) Uint(minValue, maxValue uint) (uint, error) {
	m.lastMin = minValue
	m.lastMax = maxValue

	return maxValue - 1, nil
}

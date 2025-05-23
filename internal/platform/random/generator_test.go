package random_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/platform/random"
	"testing"
)

type args struct {
	min uint
	max uint
}

type cryptoRandomTestTable struct {
	name string
	args args
}

func TestCryptoRandomGenerator(t *testing.T) {
	suite.Run(t, new(CryptoRandomGeneratorFixture))
}

type CryptoRandomGeneratorFixture struct {
	suite.Suite
	tests []cryptoRandomTestTable
}

func (crg *CryptoRandomGeneratorFixture) SetupTest() {
	for i := 0; i < 1000; i++ {
		crg.tests = append(crg.tests, cryptoRandomTestTable{name: "Test", args: args{1, 100}})
	}
}

func (crg *CryptoRandomGeneratorFixture) TestUint() {
	for _, testRow := range crg.tests {
		c := &random.Crypto{}
		got, err := c.Uint(testRow.args.min, testRow.args.max)
		crg.Require().NoError(err)
		crg.GreaterOrEqual(got, testRow.args.min)
		crg.LessOrEqual(got, testRow.args.max)
	}
}

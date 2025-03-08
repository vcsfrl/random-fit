package generator

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestGenerateSuite(t *testing.T) {
	suite.Run(t, new(GenerateSuite))
}

type GenerateSuite struct {
	suite.Suite

	generator *Generator
}

func (suite *GenerateSuite) SetupTest() {
	suite.generator = &Generator{}
}

func (suite *GenerateSuite) TestGenerateElement() {
	suite.NotNil(12)
}

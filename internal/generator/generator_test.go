package generator

import (
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/core"
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

	// Given
	definition := core.ElementDefinition{
		Metadata: core.DefinitionMetadata{
			ID:          "element-1",
			Name:        "Element 1",
			Description: "Element 1 description",
		},
		UniquePicks:  true,
		NrOfPicks:    1,
		PickStrategy: core.PickStrategyRandom,
		Options: core.ElementDefinitionOptions{
			Values: []any{"value-1", "value-2", "value-3", "value-4"},
		},
	}

	// When
	element := suite.generator.Element(definition)

	// Then
	suite.NotNil(element)
}

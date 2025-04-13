package combination

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCombinationDefinitionSuite(t *testing.T) {
	suite.Run(t, new(CombinationDefinitionSuite))
}

type CombinationDefinitionSuite struct {
	suite.Suite
}

func (suite *CombinationDefinitionSuite) TestNewCombinationDefinition() {
	script := "./testdata/star_script.star"

	definition, err := NewCombinationDefinition(script)
	suite.NoError(err)
	suite.NotNil(definition)

	suite.NotEmpty(definition.ID)
	suite.Equal("Lotto Number Picks", definition.Name)
	suite.Equal(script, definition.StarScript)
	suite.Contains(definition.GoTemplate, "/*gotype:")

	suite.NotNil(definition.buildFunction)
	combinationGenerator, err := definition.Generator()
	suite.NoError(err)

	combination, err := combinationGenerator()
	suite.NoError(err)
	suite.NotNil(combination)
}

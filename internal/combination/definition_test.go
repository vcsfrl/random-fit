package combination

import (
	"github.com/stretchr/testify/suite"
	"go.starlark.net/starlark"
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

	suite.Equal(script, definition.StarScript)
	suite.Equal("Lotto Number Picks", definition.Name)

	suite.NotNil(definition.BuildFunction)
	combinationBuilder, err := definition.CombinationBuilder()
	suite.NoError(err)

	combination, err := combinationBuilder()
	suite.NoError(err)
	suite.NotNil(combination)
	suite.Contains(combination.(*starlark.Dict).String(), "Lotto number picks")

	suite.Contains(definition.GoTemplate, "/*gotype:")
}

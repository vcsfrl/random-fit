package combination

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestCombinationDefinitionSuite(t *testing.T) {
	suite.Run(t, new(CombinationDefinitionSuite))
}

type CombinationDefinitionSuite struct {
	suite.Suite

	definition *StarlarkDefinition
	id         int
	testRand   uint
	scriptFile string
}

func (suite *CombinationDefinitionSuite) SetupTest() {
	suite.scriptFile = "./testdata/star_script.star"

	var err error
	suite.definition, err = NewCombinationDefinition(suite.scriptFile)
	suite.NoError(err)
	suite.NotNil(suite.definition)

	suite.id = 0
	suite.definition.uuidFunc = func() (string, error) {
		suite.id++
		return fmt.Sprintf("00000000-0000-0000-0000-%012d", suite.id), nil
	}

	suite.definition.nowFunc = func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	suite.definition.randomUintFunc = func(min uint, max uint) (uint, error) {
		suite.testRand++
		return suite.testRand, nil
	}

	suite.NoError(err)
}

func (suite *CombinationDefinitionSuite) TestCombinationDefinition_Generator() {
	suite.NotEmpty(suite.definition.ID)
	suite.Equal("Lotto Number Picks", suite.definition.Name)
	suite.Equal(suite.scriptFile, suite.definition.StarScript)
	suite.Contains(suite.definition.GoTemplate, "/*gotype:")

	suite.NotNil(suite.definition.buildFunction)
	combinationGenerator, err := suite.definition.Generator()
	suite.NoError(err)

	combination, err := combinationGenerator()
	suite.NoError(err)
	suite.NotNil(combination)

	suite.NotEmpty(combination.UUID)
	suite.Equal(suite.definition.ID, combination.DefinitionID)
	suite.Equal(suite.definition.Name, combination.Name)
	suite.Equal(suite.definition.GoTemplate, combination.GoTemplate)
	suite.NotEmpty(combination.Data)

	stringData := fmt.Sprintf("%+v", combination.Data)

	suite.Contains(stringData, "Name:6/49 and Lucky Number")
	suite.Contains(stringData, "User 1 Monthly Lotto Number picks")
	suite.Contains(stringData, "User 2 Monthly Lotto Number picks")
	suite.Contains(stringData, "[1 2 3 4 5 6]")
	suite.Contains(stringData, "[36 37 38 39 40 41]")
	suite.Contains(stringData, "collection_00000000-0000-0000-0000-000000000001")
	suite.Contains(stringData, "element_00000000-0000-0000-0000-000000000021")
	suite.Contains(stringData, "Lucky Number")
	suite.Contains(stringData, "Values:4200")

	suite.Equal(4205, len(stringData))
}

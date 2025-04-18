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

func (suite *CombinationDefinitionSuite) TestCombinationDefinition_CallScriptBuild() {
	buildData, err := suite.definition.CallScriptBuildFunction()
	suite.NoError(err)

	suite.Equal("lotto-test", suite.definition.ID)
	suite.Equal("Lotto Number Picks", suite.definition.Name)
	suite.Equal(suite.definition.StarScript, suite.scriptFile)

	suite.Contains(buildData, "6/49 and Lucky Number")
	suite.Contains(buildData, "Lotto Numbers for User 1")
	suite.Contains(buildData, "Lotto Numbers for User 2")
	suite.Contains(buildData, "[1,2,3,4,5,6]")
	suite.Contains(buildData, "[ 1 2 3 4 5 6 ]")
	suite.Contains(buildData, "[36,37,38,39,40,41]")
	suite.Contains(buildData, "[ 36 37 38 39 40 41 ]")
	suite.Contains(buildData, "collection_00000000-0000-0000-0000-000000000001")
	suite.Contains(buildData, "element_00000000-0000-0000-0000-000000000021")
	suite.Contains(buildData, "Lucky Number")
	suite.Contains(buildData, "4200")
}

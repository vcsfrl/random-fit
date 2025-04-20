package combination

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestStarlarkBuilder(t *testing.T) {
	suite.Run(t, new(StarlarkBuilderSuite))
}

type StarlarkBuilderSuite struct {
	suite.Suite

	definition *StarlarkDefinition
	id         int
	testRand   uint
	scriptFile string
}

func (suite *StarlarkBuilderSuite) SetupTest() {
	suite.scriptFile = "./testdata/star_script.star"
}

func (suite *StarlarkBuilderSuite) initDefinition(scriptFile string) {
	var err error
	suite.definition, err = NewCombinationDefinition(scriptFile)
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

	suite.testRand = 0
	suite.definition.randomUintFunc = func(min uint, max uint) (uint, error) {
		suite.testRand++
		return suite.testRand, nil
	}

	suite.NoError(err)
}

func (suite *StarlarkBuilderSuite) TestStarlarkBuilder_Build() {
	suite.initDefinition(suite.scriptFile)
	builder := NewStarlarkBuilder(suite.definition)
	suite.NotNil(builder)

	// Build first combination
	combination, err := builder.Build()
	suite.NoError(err)
	suite.NotNil(combination)

	suite.Equal(36, len(combination.UUID.String()))
	suite.NotNil(combination.CreatedAt)
	suite.Equal("lotto-test", combination.DefinitionID)
	suite.Equal("Lotto Number Picks", combination.Details)
	suite.NotNil(combination.Data)

	suite.Contains(combination.Data[DataTypeJson].Data.String(), "6/49 and Lucky Number")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Lotto Numbers for User 1")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Lotto Numbers for User 2")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "[1,2,3,4,5,6]")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "[36,37,38,39,40,41]")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "collection_00000000-0000-0000-0000-000000000001")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "element_00000000-0000-0000-0000-000000000021")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Lucky Number")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "4200")

	suite.Contains(combination.Data[DataTypeMd].Data.String(), "6/49 and Lucky Number")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Lotto Numbers for User 1")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Lotto Numbers for User 2")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "[ 1 2 3 4 5 6 ]")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "[ 36 37 38 39 40 41 ]")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Lucky Number")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "4200")

	// Build first combination
	combination, err = builder.Build()
	suite.NoError(err)
	suite.NotNil(combination)

	suite.Equal(36, len(combination.UUID.String()))
	suite.NotNil(combination.CreatedAt)
	suite.Equal("lotto-test", combination.DefinitionID)
	suite.Equal("Lotto Number Picks", combination.Details)
	suite.NotNil(combination.Data)

	suite.Contains(combination.Data[DataTypeJson].Data.String(), "6/49 and Lucky Number")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Lotto Numbers for User 1")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Lotto Numbers for User 2")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "[43,44,45,46,47,48]")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "[78,79,80,81,82,83]")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "collection_00000000-0000-0000-0000-000000000022")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "set_00000000-0000-0000-0000-000000000040")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Lucky Number")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "8400")

	suite.Contains(combination.Data[DataTypeMd].Data.String(), "6/49 and Lucky Number")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Lotto Numbers for User 1")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Lotto Numbers for User 2")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "[ 43 44 45 46 47 48 ]")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "[ 78 79 80 81 82 83 ]")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Lucky Number")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "8400")
}

func (suite *StarlarkBuilderSuite) TestStarlarkBuilder_NoJsonData() {
	suite.initDefinition("./testdata/star_script_no_json.star")
	builder := NewStarlarkBuilder(suite.definition)
	suite.NotNil(builder)

	// Build first combination
	combination, err := builder.Build()
	suite.Error(err, "combination data does not contain json representation (required)")
	suite.Nil(combination)
}

func (suite *StarlarkBuilderSuite) TestStarlarkBuilder_Sample() {
	suite.initDefinition("./template/script.star")
	builder := NewStarlarkBuilder(suite.definition)
	suite.NotNil(builder)

	// Build first combination
	combination, err := builder.Build()
	suite.Nil(err)
	suite.NotNil(combination)

	suite.Equal(36, len(combination.UUID.String()))
	suite.NotNil(combination.CreatedAt)
	suite.Equal("sample", combination.DefinitionID)
	suite.Equal("Sample Combination", combination.Details)
	suite.NotNil(combination.Data)
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "Sample")
	suite.Contains(combination.Data[DataTypeJson].Data.String(), "[1,2,3,4,5,6,7,8,9,10]")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "Sample")
	suite.Contains(combination.Data[DataTypeMd].Data.String(), "[ 1 2 3 4 5 6 7 8 9 10 ]")
}

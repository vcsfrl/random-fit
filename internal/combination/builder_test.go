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

func (suite *StarlarkBuilderSuite) TestStarlarkBuilder_Build() {
	builder := NewStarlarkBuilder(suite.definition)
	suite.NotNil(builder)

	// Build first combination
	combination, err := builder.Build()
	suite.NoError(err)
	suite.NotNil(combination)

	suite.Equal(36, len(combination.UUID.String()))
	suite.NotNil(combination.CreatedAt)
	suite.Equal("lotto-test", combination.DefinitionID)
	suite.Equal("Lotto Number Picks", combination.Name)
	suite.Contains(combination.Template, "{{- /*Generate lotto numbers*/ -}}")
	suite.NotNil(combination.Data)

	suite.Contains(combination.Data.String(), "6/49 and Lucky Number")
	suite.Contains(combination.Data.String(), "Lotto Numbers for User 1")
	suite.Contains(combination.Data.String(), "Lotto Numbers for User 2")
	suite.Contains(combination.Data.String(), "[1, 2, 3, 4, 5, 6]")
	suite.Contains(combination.Data.String(), "[36, 37, 38, 39, 40, 41]")
	suite.Contains(combination.Data.String(), "collection_00000000-0000-0000-0000-000000000001")
	suite.Contains(combination.Data.String(), "element_00000000-0000-0000-0000-000000000021")
	suite.Contains(combination.Data.String(), "Lucky Number")
	suite.Contains(combination.Data.String(), "4200")
	suite.Equal(4354, len(combination.Data.String()))

	// Build first combination
	combination, err = builder.Build()
	suite.NoError(err)
	suite.NotNil(combination)

	suite.Equal(36, len(combination.UUID.String()))
	suite.NotNil(combination.CreatedAt)
	suite.Equal("lotto-test", combination.DefinitionID)
	suite.Equal("Lotto Number Picks", combination.Name)
	suite.Contains(combination.Template, "{{- /*Generate lotto numbers*/ -}}")
	suite.NotNil(combination.Data)

	suite.Contains(combination.Data.String(), "6/49 and Lucky Number")
	suite.Contains(combination.Data.String(), "Lotto Numbers for User 1")
	suite.Contains(combination.Data.String(), "Lotto Numbers for User 2")
	suite.Contains(combination.Data.String(), "[43, 44, 45, 46, 47, 48]")
	suite.Contains(combination.Data.String(), "[78, 79, 80, 81, 82, 83]")
	suite.Contains(combination.Data.String(), "collection_00000000-0000-0000-0000-000000000022")
	suite.Contains(combination.Data.String(), "set_00000000-0000-0000-0000-000000000040")
	suite.Contains(combination.Data.String(), "Lucky Number")
	suite.Contains(combination.Data.String(), "8400")
	suite.Equal(4363, len(combination.Data.String()))

}

//
//func (suite *StarlarkBuilderSuite) TestStarlarkBuilder_Build_View() {
//
//	builder := NewStarlarkBuilder(suite.definition)
//	suite.NotNil(builder)
//
//	// Build first combination
//	combination, err := builder.Build()
//
//	tmpl, err := template.New(combination.DefinitionID).Parse(combination.Template)
//	suite.NoError(err)
//
//	buffer := &bytes.Buffer{}
//	//buffer, err := os.OpenFile("testdata/lotto_test.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
//	suite.NoError(err)
//	err = tmpl.Execute(buffer, combination)
//	suite.NoError(err)
//
//	fmt.Println(buffer.String())
//
//}

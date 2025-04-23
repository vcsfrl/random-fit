package plan

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/combination"
	"os"
	"path/filepath"
	"testing"
)

func TestExportSuite(t *testing.T) {
	suite.Run(t, new(ExportSuite))
}

type ExportSuite struct {
	suite.Suite

	testFolder         string
	planDefinition     *Definition
	combinationBuilder combination.Builder
	planBuilder        *Builder
}

func (suite *ExportSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	// Create a test plan definition
	suite.planDefinition = &Definition{
		ID:      "test",
		Details: "Test",
		Users:   []string{"user-1"},
		GroupDefinition: GroupDefinition{
			NamePrefix:       "Test",
			NumberOfGroups:   4,
			NrOfCombinations: 3,
		},
	}

	definition, err := combination.NewCombinationDefinition("./testdata/star_script.star")
	suite.NoError(err)
	suite.combinationBuilder = combination.NewStarlarkBuilder(definition)

	suite.planBuilder = NewBuilder(suite.planDefinition, suite.combinationBuilder)
}

func (suite *ExportSuite) TearDownTest() {
	// Remove the test folder
	//err := os.RemoveAll(suite.testFolder)
	//suite.NoError(err)
}

func (suite *ExportSuite) TestExport() {
	//plan, err := suite.planBuilder.Build()
	//suite.NoError(err)
	//suite.NotNil(plan)

	//for _, user := range plan.Users {
	//	for _, group := range plan.UserGroups {
	//		fmt.Println(user, group)
	//	}
	//}

}

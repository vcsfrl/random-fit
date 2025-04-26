package plan

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/combination"
	"github.com/vcsfrl/random-fit/internal/platform/starlark/random"
	slUuid "github.com/vcsfrl/random-fit/internal/platform/starlark/uuid"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	id                 int
	testRand           uint
}

func (suite *ExportSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	// Create a test plan definition
	suite.planDefinition = &Definition{
		ID:      "test-definition",
		Details: "Test definition",
		Users:   []string{"user-1"},
		UserData: UserData{
			ContainerName:            "Group-Container",
			RecurrentGroupNamePrefix: "Recurrent-Group",
			RecurrentGroups:          4,
			NrOfGroupCombinations:    3,
		},
	}

	definition, err := combination.NewCombinationDefinition("./testdata/star_script.star")
	suite.NoError(err)
	suite.combinationBuilder = combination.NewStarlarkBuilder(definition)

	suite.id = 0
	slUuid.SetUuidFunc(func() (string, error) {
		suite.id++
		return fmt.Sprintf("00000000-0000-0000-0000-%012d", suite.id), nil
	})

	suite.testRand = 0
	random.SetUintFunc(func(min uint, max uint) (uint, error) {
		suite.testRand++
		return suite.testRand, nil
	})

	suite.planBuilder = NewBuilder(suite.planDefinition, suite.combinationBuilder)
	suite.planBuilder.Now = func() time.Time {
		return time.Date(2010, 1, 2, 3, 4, 5, 6, time.UTC)
	}
}

func (suite *ExportSuite) TearDownTest() {
	// Remove the test folder
	err := os.RemoveAll(suite.testFolder)
	suite.NoError(err)
}

func (suite *ExportSuite) TestExport() {
	plan, err := suite.planBuilder.Build()
	suite.NoError(err)
	suite.NotNil(plan)

	exporter := NewExporter(suite.testFolder)
	err = exporter.Export(plan)
}

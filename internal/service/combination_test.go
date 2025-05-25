package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/platform/fs"
	"github.com/vcsfrl/random-fit/internal/service"
	"os"
	"path/filepath"
	"testing"
)

const testDefinitionFileName = "test-definitionFileName"

func TestDefinitionManager(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StarDefinitionManagerSuite))
}

type StarDefinitionManagerSuite struct {
	suite.Suite
	testFolder string

	definitionManager *service.CombinationStarDefinitionManager
}

func (suite *StarDefinitionManagerSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.Require().NoError(err)

	suite.definitionManager = service.NewCombinationStarDefinitionManager(suite.testFolder)
}

func (suite *StarDefinitionManagerSuite) TearDownTest() {
	// Remove the test folder
	err := os.RemoveAll(suite.testFolder)
	suite.Require().NoError(err)
}

func (suite *StarDefinitionManagerSuite) TestList() {
	// create a test definitionFileName files
	testDefinitions := []string{"test-definitionFileName-1", "test-definitionFileName-2", "test-definitionFileName-3"}
	for _, definitionFileName := range testDefinitions {
		testDefinitionFile := filepath.Join(suite.testFolder, definitionFileName+".star")
		err := os.WriteFile(testDefinitionFile, []byte(`test`), fs.FilePermission)
		suite.Require().NoError(err)
	}

	definitions, err := suite.definitionManager.List()
	suite.Require().NoError(err)
	suite.NotNil(definitions)
	suite.Len(definitions, len(testDefinitions))

	for _, definitionFileName := range testDefinitions {
		suite.Contains(definitions, definitionFileName)
	}
}

func (suite *StarDefinitionManagerSuite) TestNew() {
	// create a test definitionFileName file
	testDefinitionFile := filepath.Join(suite.testFolder, testDefinitionFileName+".star")

	err := suite.definitionManager.New(testDefinitionFileName)
	suite.Require().NoError(err)

	// check if the file exists
	_, err = os.Stat(testDefinitionFile)
	suite.Require().NoError(err)

	data, err := os.ReadFile(testDefinitionFile)
	suite.Require().NoError(err)
	suite.NotEmpty(data)
	suite.Equal(service.DefinitionTemplate, string(data))

	// do not overwrite the file if it already exists
	err = suite.definitionManager.New(testDefinitionFileName)
	suite.Require().Error(err)
}

func (suite *StarDefinitionManagerSuite) TestGetScript() {
	// create a test definitionFileName file
	testDefinitionFile := filepath.Join(suite.testFolder, testDefinitionFileName+".star")

	err := suite.definitionManager.New(testDefinitionFileName)
	suite.Require().NoError(err)

	script, err := suite.definitionManager.GetScript(testDefinitionFileName)
	suite.Require().NoError(err)
	suite.NotEmpty(script)
	suite.Equal(testDefinitionFile, script)
}

func (suite *StarDefinitionManagerSuite) TestBuild() {
	err := suite.definitionManager.New(testDefinitionFileName)
	suite.Require().NoError(err)

	combination, err := suite.definitionManager.Build(testDefinitionFileName)
	suite.Require().NoError(err)
	suite.NotNil(combination)

	suite.Equal("Sample Combination", combination.Details)
}

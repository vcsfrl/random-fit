package internal

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

func TestDefinitionManager(t *testing.T) {
	suite.Run(t, new(StarDefinitionManagerSuite))
}

type StarDefinitionManagerSuite struct {
	suite.Suite
	testFolder string

	definitionManager *CombinationStarDefinitionManager
}

func (suite *StarDefinitionManagerSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	suite.definitionManager = NewCombinationStarDefinitionManager(suite.testFolder)
}

func (suite *StarDefinitionManagerSuite) TearDownTest() {
	// Remove the test folder
	err := os.RemoveAll(suite.testFolder)
	suite.NoError(err)
}

func (suite *StarDefinitionManagerSuite) TestList() {
	// create a test definitionFileName files
	testDefinitions := []string{"test-definitionFileName-1", "test-definitionFileName-2", "test-definitionFileName-3"}
	for _, definitionFileName := range testDefinitions {
		testDefinitionFile := filepath.Join(suite.testFolder, fmt.Sprintf("%s.star", definitionFileName))
		err := os.WriteFile(testDefinitionFile, []byte(`test`), 0644)
		suite.NoError(err)
	}

	definitions, err := suite.definitionManager.List()
	suite.NoError(err)
	suite.NotNil(definitions)
	suite.Equal(len(testDefinitions), len(definitions))

	for _, definitionFileName := range testDefinitions {
		suite.Contains(definitions, definitionFileName)
	}
}

func (suite *StarDefinitionManagerSuite) TestNew() {
	// create a test definitionFileName file
	testDefinitionFileName := "test-definitionFileName"
	testDefinitionFile := filepath.Join(suite.testFolder, fmt.Sprintf("%s.star", testDefinitionFileName))

	err := suite.definitionManager.New(testDefinitionFileName)
	suite.NoError(err)

	// check if the file exists
	_, err = os.Stat(testDefinitionFile)
	suite.NoError(err)

	data, err := os.ReadFile(testDefinitionFile)
	suite.NoError(err)
	suite.NotEmpty(data)
	suite.Equal(definitionTemplate, string(data))

	// do not overwrite the file if it already exists
	err = suite.definitionManager.New(testDefinitionFileName)
	suite.Error(err)
}

func (suite *StarDefinitionManagerSuite) TestGetScript() {
	// create a test definitionFileName file
	testDefinitionFileName := "test-definitionFileName"
	testDefinitionFile := filepath.Join(suite.testFolder, fmt.Sprintf("%s.star", testDefinitionFileName))

	err := suite.definitionManager.New(testDefinitionFileName)
	suite.NoError(err)

	script, err := suite.definitionManager.GetScript(testDefinitionFileName)
	suite.NoError(err)
	suite.NotEmpty(script)
	suite.Equal(testDefinitionFile, script)
}

func (suite *StarDefinitionManagerSuite) TestBuild() {
	testDefinitionFileName := "test-definitionFileName"
	err := suite.definitionManager.New(testDefinitionFileName)
	suite.NoError(err)

	combination, err := suite.definitionManager.Build(testDefinitionFileName)
	suite.NoError(err)
	suite.NotNil(combination)

	suite.Equal("Sample Combination", combination.Details)
}

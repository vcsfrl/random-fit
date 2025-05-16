package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	rfPlan "github.com/vcsfrl/random-fit/internal/plan"
	"os"
	"path/filepath"
	"testing"
)

func TestPlanDefinitionManager(t *testing.T) {
	suite.Run(t, new(StarPlanManagerSuite))
}

type StarPlanManagerSuite struct {
	suite.Suite
	testFolder string

	planDefinitionManager *PlanDefinitionManager
}

func (suite *StarPlanManagerSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	suite.planDefinitionManager = NewPlanDefinitionManager(suite.testFolder)
}

func (suite *StarPlanManagerSuite) TearDownTest() {
	// Remove the test folder
	err := os.RemoveAll(suite.testFolder)
	suite.NoError(err)
}

func (suite *StarPlanManagerSuite) TestList() {
	// create a test plan files
	testPlans := []string{"test-plan-1", "test-plan-2", "test-plan-3"}
	for _, plan := range testPlans {
		testPlanFile := filepath.Join(suite.testFolder, plan)
		err := os.WriteFile(testPlanFile, []byte(`test`), 0644)
		suite.NoError(err)
	}

	plans, err := suite.planDefinitionManager.List()
	suite.NoError(err)
	suite.NotNil(plans)
	suite.Equal(len(testPlans), len(plans))

	for _, plan := range testPlans {
		suite.Contains(plans, plan)
	}
}

func (suite *StarPlanManagerSuite) TestNew() {
	testPlan := "test-plan"
	// create a test plan file
	testPlanFile := filepath.Join(suite.testFolder, fmt.Sprintf("%s.json", testPlan))

	// create a new plan
	err := suite.planDefinitionManager.New(testPlan)
	suite.NoError(err)

	// check if the plan file exists
	_, err = os.Stat(testPlanFile)
	suite.NoError(err)

	// check if the plan file is empty
	fileInfo, err := os.Stat(testPlanFile)
	suite.NoError(err)
	suite.Greater(fileInfo.Size(), int64(0))

	// check if the plan file is valid json
	data, err := os.ReadFile(testPlanFile)
	suite.NoError(err)
	suite.NotEmpty(data)

	resultPlanDefinition := &rfPlan.Definition{}

	err = json.Unmarshal(data, resultPlanDefinition)
	suite.NoError(err)
	suite.Equal(suite.planDefinitionManager.getSamplePlanDefinition(), resultPlanDefinition)
}

func (suite *StarPlanManagerSuite) TestGetFile() {
	testPlan := "test-plan"
	// create a test plan file
	testPlanFile := filepath.Join(suite.testFolder, fmt.Sprintf("%s.json", testPlan))

	// create a new plan
	err := suite.planDefinitionManager.New(testPlan)
	suite.NoError(err)

	// check if the plan file exists
	_, err = os.Stat(testPlanFile)
	suite.NoError(err)

	// get the plan file
	result, err := suite.planDefinitionManager.GetFile(testPlan)
	suite.NoError(err)
	suite.Equal(testPlanFile, result)
}

package shell

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

func TestPlanManager(t *testing.T) {
	suite.Run(t, new(StarPlanManagerSuite))
}

type StarPlanManagerSuite struct {
	suite.Suite
	testFolder string

	planManager *PlanManager
}

func (suite *StarPlanManagerSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	suite.planManager = NewPlanManager(suite.testFolder)
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

	plans, err := suite.planManager.List()
	suite.NoError(err)
	suite.NotNil(plans)
	suite.Equal(len(testPlans), len(plans))

	for _, plan := range testPlans {
		suite.Contains(plans, plan)
	}
}

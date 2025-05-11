package shell

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestShell(t *testing.T) {
	suite.Run(t, new(ShellSuite))
}

type ShellSuite struct {
	suite.Suite

	shell *Shell
}

func (suite *ShellSuite) SetupTest() {
	// Create a new shell instance
	suite.shell = New()

	// Check if the shell instance is not nil
	suite.NotNil(suite.shell)

	suite.NotNil(suite.shell.stdin)
	suite.NotNil(suite.shell.stdout)
	suite.NotNil(suite.shell.stderr)
}

func (suite *ShellSuite) TearDownTest() {
	// Close the shell instance
	if suite.shell != nil {
		_ = suite.shell.Close()
	}
}

func (suite *ShellSuite) TestNew() {
	// Create a new shell instance
	newShell := New()

	// Check if the shell instance is not nil
	suite.NotNil(newShell)

	// Check if the shell instance has the expected properties

	suite.NotNil(newShell.definitionFolder)
	suite.NotNil(newShell.planFolder)
	suite.NotNil(newShell.storageFolder)
	suite.NotNil(newShell.combinationFolder)
}

package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCommands(t *testing.T) {
	suite.Run(t, new(CommandsSuite))
}

type CommandsSuite struct {
	suite.Suite

	buffer  *bytes.Buffer
	command *cobra.Command
}

func (suite *CommandsSuite) SetupTest() {
	// Capture output
	var err error
	suite.buffer = new(bytes.Buffer)
	suite.command, err = NewCommand()
	suite.NoError(err)

	suite.command.SetOut(suite.buffer)
	suite.command.SetErr(suite.buffer)
}

func (suite *CommandsSuite) TearDownTest() {

}

func (suite *CommandsSuite) TestSubcommands() {
	// Assert subcommands are present
	subcommands := suite.command.Commands()
	var subcommandNames []string
	for _, cmd := range subcommands {
		subcommandNames = append(subcommandNames, cmd.Name())
	}

	suite.Len(subcommandNames, 2)
	suite.Contains(subcommandNames, "definition")
	suite.Contains(subcommandNames, "code")
}

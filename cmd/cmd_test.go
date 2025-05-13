package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"slices"
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
	subcommands := suite.command.Commands()
	var expectedSubcommandNames = []string{"definition", "code", "generate"}
	var subcommandNames []string
	for _, cmd := range subcommands {
		subcommandNames = append(subcommandNames, cmd.Name())
		for _, subCmd := range cmd.Commands() {
			subcommandNames = append(subcommandNames, subCmd.Name())
		}

	}

	slices.Sort(subcommandNames)
	slices.Sort(expectedSubcommandNames)

	suite.Equal(expectedSubcommandNames, subcommandNames)
}

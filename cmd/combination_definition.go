package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/cmd/internal"
	"os"
)

type CombinationDefinition struct {
	BaseDefinition
}

func NewCombinationDefinition(cmd *cobra.Command, args []string, conf *internal.Config) *CombinationDefinition {
	return &CombinationDefinition{
		BaseDefinition: BaseDefinition{
			cmd:  cmd,
			args: args,
			conf: conf,
		},
	}
}

func (c *CombinationDefinition) New() {
	name := c.getNameArg()
	if name == "" {
		c.cmd.PrintErrln(msgNameMissing)
		return
	}

	c.cmd.Println(msqCreate, msgCombinationDefinition, name)
	if err := createFolder(c.conf.DefinitionFolder()); err != nil {
		c.cmd.PrintErrln("Error creating definition folder: ", err)
		return
	}

	definitionManager := internal.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())
	err := definitionManager.New(name)
	if err != nil {
		c.cmd.PrintErrln("Error: ", err)
		return
	}

	c.cmd.Println(msgCombinationDefinition, msqCreated, name)
	scriptName, err := definitionManager.GetScript(name)
	if err != nil {
		c.cmd.PrintErrln("Error getting script: ", err)
		return
	}

	c.cmd.Println(msqEditScript, scriptName)
	if err := c.editScript(scriptName, "python"); err != nil {
		c.cmd.PrintErrln("Error editing script: ", err)
		return
	}
}

func createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}

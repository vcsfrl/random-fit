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

	c.cmd.Println(msgCreate, msgCombinationDefinition, name)
	if err := c.createFolder(c.conf.DefinitionFolder()); err != nil {
		c.cmd.PrintErrln("Error creating definition folder: ", err)
		return
	}

	definitionManager := internal.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())
	err := definitionManager.New(name)
	if err != nil {
		c.cmd.PrintErrln("Error: ", err)
		return
	}

	c.cmd.Println(msgDone, msgCreate, msgCombinationDefinition, name)
	scriptName, err := definitionManager.GetScript(name)
	if err != nil {
		c.cmd.PrintErrln("Error getting script: ", err)
		return
	}

	c.cmd.Println(msgEditScript, scriptName)
	if err := c.editScript(scriptName, "python"); err != nil {
		c.cmd.PrintErrln("Error editing script: ", err)
		return
	}
}

func (c *CombinationDefinition) Edit() {
	name := c.getNameArg()
	if name == "" {
		c.cmd.PrintErrln(msgNameMissing)
		return
	}

	c.cmd.Println(msgEdit, msgCombinationDefinition, name)
	if err := c.createFolder(c.conf.DefinitionFolder()); err != nil {
		c.cmd.PrintErrln("Error creating definition folder: ", err)
		return
	}

	definitionManager := internal.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())
	scriptName, err := definitionManager.GetScript(name)
	if err != nil {
		c.cmd.PrintErrln("Error getting script: ", err)
		return
	}

	c.cmd.Println(msgEditScript, scriptName)
	if err := c.editScript(scriptName, "python"); err != nil {
		c.cmd.PrintErrln("Error editing script: ", err)
		return
	}

}

func (c *CombinationDefinition) Delete() {
	name := c.getNameArg()
	if name == "" {
		c.cmd.PrintErrln(msgNameMissing)
		return
	}

	c.cmd.Println(msgDelete, msgCombinationDefinition, name)
	if err := c.createFolder(c.conf.DefinitionFolder()); err != nil {
		c.cmd.PrintErrln("Error creating definition folder: ", err)
		return
	}

	definitionManager := internal.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())
	scriptName, err := definitionManager.GetScript(name)
	if err != nil {
		c.cmd.PrintErrln("Error getting script: ", err)
		return
	}

	c.cmd.Println(msgRemoveScript, scriptName)
	if err := os.Remove(scriptName); err != nil {
		c.cmd.PrintErrln("Error removing script: ", err)
		return
	}
}

func (c *CombinationDefinition) List() {
	c.cmd.Println(msgCombinationDefinition, msgList)
	if err := c.createFolder(c.conf.DefinitionFolder()); err != nil {
		c.cmd.PrintErrln("Error creating definition folder: ", err)
		return
	}

	definitionManager := internal.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())
	definitions, err := definitionManager.List()
	if err != nil {
		c.cmd.PrintErrln("Error listing definitions: ", err)
		return
	}

	if len(definitions) == 0 {
		c.cmd.Println(msgNoItemsFound)
		return
	}

	for _, definition := range definitions {
		c.cmd.Println(" - " + definition)
	}
}

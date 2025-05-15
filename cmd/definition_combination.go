package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/cmd/internal"
	"os"
)

type CombinationDefinition struct {
	BaseHandler

	definitionManager *internal.CombinationStarDefinitionManager
}

func NewCombinationDefinition(cmd *cobra.Command, args []string, conf *internal.Config) (*CombinationDefinition, error) {
	combinationDefinition := &CombinationDefinition{
		BaseHandler: BaseHandler{
			cmd:  cmd,
			args: args,
			conf: conf,
		},
	}

	if err := combinationDefinition.init(); err != nil {
		return nil, err
	}

	return combinationDefinition, nil
}

func (c *CombinationDefinition) init() error {
	err := c.initFolders()
	if err != nil {
		return err
	}

	c.definitionManager = internal.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())
	return nil
}

func (c *CombinationDefinition) New() {
	name := c.getArg(0, "name")
	if name == "" {
		c.cmd.PrintErrln(msgNameMissing)
		return
	}

	c.cmd.Println(msgCreate, msgCombinationDefinition, name)
	err := c.definitionManager.New(name)
	if err != nil {
		c.cmd.PrintErrln("Error: ", err)
		return
	}

	c.cmd.Println(msgDone, msgCreate, msgCombinationDefinition, name)
	scriptName, err := c.definitionManager.GetScript(name)
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
	name := c.getArg(0, "name")
	if name == "" {
		c.cmd.PrintErrln(msgNameMissing)
		return
	}

	c.cmd.Println(msgEdit, msgCombinationDefinition, name)
	scriptName, err := c.definitionManager.GetScript(name)
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
	name := c.getArg(0, "name")
	if name == "" {
		c.cmd.PrintErrln(msgNameMissing)
		return
	}

	c.cmd.Println(msgDelete, msgCombinationDefinition, name)
	scriptName, err := c.definitionManager.GetScript(name)
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
	definitions, err := c.definitionManager.List()
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

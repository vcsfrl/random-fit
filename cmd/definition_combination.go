package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/service"
)

type CombinationDefinition struct {
	BaseHandler

	definitionManager *service.CombinationStarDefinitionManager
}

func NewCombinationDefinition(cmd *cobra.Command, args []string, conf *service.Config) (*CombinationDefinition, error) {
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

func (c *CombinationDefinition) New() {
	name := c.getArg(0, "name")
	if name == "" {
		c.cmd.PrintErrln(c.printer.Sprint("Name is required."))

		return
	}

	c.cmd.Println(c.printer.Sprint("Create combination definition:"), name)
	err := c.definitionManager.New(name)

	if err != nil {
		c.cmd.PrintErrln("Error: ", err)

		return
	}

	c.cmd.Println(c.printer.Sprint("Finished creating combination definition:"), name)
	scriptName, err := c.definitionManager.GetScript(name)

	if err != nil {
		c.cmd.PrintErrln("Error getting script: ", err)

		return
	}

	c.cmd.Println(c.printer.Sprint("Edit combination definition file:"), scriptName)

	if err := c.editScript(scriptName, "python"); err != nil {
		c.cmd.PrintErrln("Error editing script: ", err)

		return
	}
}

func (c *CombinationDefinition) Edit() {
	name := c.getArg(0, "name")
	if name == "" {
		c.cmd.PrintErrln(c.printer.Sprint("Name is required."))

		return
	}

	c.cmd.Println(c.printer.Sprint("Edit combination definition: "), name)

	scriptName, err := c.definitionManager.GetScript(name)
	if err != nil {
		c.cmd.PrintErrln("Error getting script: ", err)

		return
	}

	c.cmd.Println(c.printer.Sprint("Edit script:"), scriptName)

	if err := c.editScript(scriptName, "python"); err != nil {
		c.cmd.PrintErrln("Error editing script: ", err)

		return
	}
}

func (c *CombinationDefinition) Delete() {
	name := c.getArg(0, "name")
	if name == "" {
		c.cmd.PrintErrln(c.printer.Sprint("Name is required."))

		return
	}

	c.cmd.Println(c.printer.Sprint("Delete combination definition:"), name)
	err := c.definitionManager.Delete(name)

	if err != nil {
		c.cmd.PrintErrln("Error deleting script: ", err)

		return
	}

	c.cmd.Println(c.printer.Sprint("Finished deleting combination definition:"), name)
}

func (c *CombinationDefinition) List() {
	c.cmd.Println(c.printer.Sprint("Combination definitions:"))
	definitions, err := c.definitionManager.List()

	if err != nil {
		c.cmd.PrintErrln("Error listing definitions: ", err)

		return
	}

	if len(definitions) == 0 {
		c.cmd.Println(c.printer.Sprint("No combination definitions found."))

		return
	}

	for _, definition := range definitions {
		c.cmd.Println(" - " + definition)
	}
}

func (c *CombinationDefinition) init() error {
	c.initTranslations()

	err := c.initFolders()
	if err != nil {
		return err
	}

	c.definitionManager = service.NewCombinationStarDefinitionManager(c.conf.DefinitionFolder())

	return nil
}

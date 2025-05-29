package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/service"
)

type PlanDefinition struct {
	BaseHandler

	definitionManager *service.PlanDefinitionManager
}

func NewPlanDefinition(cmd *cobra.Command, args []string, conf *Config) (*PlanDefinition, error) {
	planDefinition := &PlanDefinition{
		BaseHandler: BaseHandler{
			cmd:  cmd,
			args: args,
			conf: conf,
		},
	}

	if err := planDefinition.init(); err != nil {
		return nil, err
	}

	return planDefinition, nil
}

func (p *PlanDefinition) New() {
	name := p.getArg(0, "name")
	if name == "" {
		p.cmd.PrintErrln(p.printer.Sprintf("Name is required."))

		return
	}

	p.cmd.Println(p.printer.Sprintf("Creating plan definition:"), name)

	err := p.definitionManager.New(name)
	if err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error:"), err)

		return
	}

	p.cmd.Println(p.printer.Sprintf("Finished creating plan definition:"), name)

	scriptName, err := p.definitionManager.GetFile(name)
	if err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error getting script:"), err)

		return
	}

	p.cmd.Println(p.printer.Sprintf("Edit definition file:"), scriptName)

	if err := p.editScript(scriptName, "json"); err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error editing script:"), err)

		return
	}
}

func (p *PlanDefinition) List() {
	p.cmd.Println(p.printer.Sprintf("Plan Definitions:"))

	planDefinitions, err := p.definitionManager.List()
	if err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error:"), err)

		return
	}

	if len(planDefinitions) == 0 {
		p.cmd.Println(p.printer.Sprintf("No plan definitions found."))

		return
	}

	for _, plan := range planDefinitions {
		p.cmd.Println(" - " + plan)
	}
}

func (p *PlanDefinition) Edit() {
	name := p.getArg(0, "name")
	if name == "" {
		p.cmd.PrintErrln(p.printer.Sprintf("Name is required."))

		return
	}

	p.cmd.Println(p.printer.Sprintf("Editing plan definition:"), name)
	scriptName, err := p.definitionManager.GetFile(name)

	if err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error getting script:"), err)

		return
	}

	if err := p.editScript(scriptName, "json"); err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error editing script:"), err)

		return
	}

	p.cmd.Println(p.printer.Sprintf("Finished editing plan definition:"), name)
}

func (p *PlanDefinition) Delete() {
	name := p.getArg(0, "name")
	if name == "" {
		p.cmd.PrintErrln(p.printer.Sprintf("Name is required."))

		return
	}

	p.cmd.Println(p.printer.Sprintf("Delete plan definition:"), name)

	err := p.definitionManager.Delete(name)
	if err != nil {
		p.cmd.PrintErrln(p.printer.Sprintf("Error deleting plan definition:"), err)

		return
	}

	p.cmd.Println(p.printer.Sprintf("Finished deleting plan definition:"), name)
}

func (p *PlanDefinition) init() error {
	p.initTranslations()

	err := p.initFolders()
	if err != nil {
		return err
	}

	p.definitionManager = service.NewPlanDefinitionManager(p.conf.PlanFolder())

	return nil
}

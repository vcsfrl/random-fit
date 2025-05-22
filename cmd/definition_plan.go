package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/service"
)

type PlanDefinition struct {
	BaseHandler

	definitionManager *service.PlanDefinitionManager
}

func NewPlanDefinition(cmd *cobra.Command, args []string, conf *service.Config) (*PlanDefinition, error) {
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
		p.cmd.PrintErrln(msgNameMissing)
		return
	}

	p.cmd.Println(msgCreate, msgPlanDefinition, name)

	err := p.definitionManager.New(name)
	if err != nil {
		p.cmd.PrintErrln("Error: ", err)
		return
	}

	p.cmd.Println(msgDone, msgCreate, msgPlanDefinition, name)

	scriptName, err := p.definitionManager.GetFile(name)
	if err != nil {
		p.cmd.PrintErrln("Error getting script: ", err)
		return
	}

	p.cmd.Println(msgEditScript, scriptName)

	if err := p.editScript(scriptName, "python"); err != nil {
		p.cmd.PrintErrln("Error editing script: ", err)
		return
	}
}

func (p *PlanDefinition) List() {
	p.cmd.Println(msgPlanDefinition, msgList)

	planDefinitions, err := p.definitionManager.List()
	if err != nil {
		p.cmd.PrintErrln("Error: ", err)
		return
	}

	if len(planDefinitions) == 0 {
		p.cmd.Println(msgNoItemsFound)
		return
	}

	for _, plan := range planDefinitions {
		p.cmd.Println(" - " + plan)
	}
}

func (p *PlanDefinition) Edit() {
	name := p.getArg(0, "name")
	if name == "" {
		p.cmd.PrintErrln(msgNameMissing)
		return
	}

	p.cmd.Println(msgEdit, msgPlanDefinition, name)
	scriptName, err := p.definitionManager.GetFile(name)

	if err != nil {
		p.cmd.PrintErrln("Error getting script: ", err)
		return
	}

	if err := p.editScript(scriptName, "python"); err != nil {
		p.cmd.PrintErrln("Error editing script: ", err)
		return
	}

	p.cmd.Println(msgDone, msgEdit, msgPlanDefinition, name)
}

func (p *PlanDefinition) Delete() {
	name := p.getArg(0, "name")
	if name == "" {
		p.cmd.PrintErrln(msgNameMissing)
		return
	}

	p.cmd.Println(msgDelete, msgPlanDefinition, name)

	err := p.definitionManager.Delete(name)
	if err != nil {
		p.cmd.PrintErrln("Error deleting paln definition: ", err)
		return
	}

	p.cmd.Println(msgDone, msgDelete, msgPlanDefinition, name)
}

func (p *PlanDefinition) init() error {
	err := p.initFolders()
	if err != nil {
		return err
	}

	p.definitionManager = service.NewPlanDefinitionManager(p.conf.PlanFolder())
	return nil
}

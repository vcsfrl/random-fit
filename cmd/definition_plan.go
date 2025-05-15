package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/cmd/internal"
)

type PlanDefinition struct {
	BaseDefinition
}

func NewPlanDefinition(cmd *cobra.Command, args []string, conf *internal.Config) *PlanDefinition {
	return &PlanDefinition{
		BaseDefinition: BaseDefinition{
			cmd:  cmd,
			args: args,
			conf: conf,
		},
	}
}

func (p *PlanDefinition) New() {
	name := p.getNameArg()
	if name == "" {
		p.cmd.PrintErrln(msgNameMissing)
		return
	}

	p.cmd.Println(msgCreate, msgPlanDefinition, name)
	if err := p.createFolder(p.conf.PlanFolder()); err != nil {
		p.cmd.PrintErrln("Error creating plan folder: ", err)
		return
	}

	definitionManager := internal.NewPlanDefinitionManager(p.conf.PlanFolder())
	err := definitionManager.New(name)
	if err != nil {
		p.cmd.PrintErrln("Error: ", err)
		return
	}

	p.cmd.Println(msgDone, msgCreate, msgPlanDefinition, name)
	scriptName, err := definitionManager.GetFile(name)
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
	if err := p.createFolder(p.conf.PlanFolder()); err != nil {
		p.cmd.PrintErrln("Error creating plan folder: ", err)
		return
	}

	definitionManager := internal.NewPlanDefinitionManager(p.conf.PlanFolder())
	planDefinitions, err := definitionManager.List()
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

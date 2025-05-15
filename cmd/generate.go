package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/cmd/internal"
)

func NewGenerator(cmd *cobra.Command, args []string, config *internal.Config) (*Generator, error) {
	generator := &Generator{
		BaseHandler: BaseHandler{
			cmd:  cmd,
			args: args,
			conf: config,
		},
	}

	generator.init()

	return generator, nil
}

type Generator struct {
	BaseHandler

	combinationDefinitionManager *internal.CombinationStarDefinitionManager
	planDefinitionManager        *internal.PlanDefinitionManager
}

func (g *Generator) Combination() {

}

func (g *Generator) init() error {
	err := g.initFolders()
	if err != nil {
		return err
	}

	g.combinationDefinitionManager = internal.NewCombinationStarDefinitionManager(g.conf.DefinitionFolder())
	g.planDefinitionManager = internal.NewPlanDefinitionManager(g.conf.PlanFolder())

	return nil
}

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/plan"
	"github.com/vcsfrl/random-fit/internal/service"
	"time"
)

func NewGenerator(cmd *cobra.Command, args []string, config *service.Config) (*Generator, error) {
	generator := &Generator{
		BaseHandler: BaseHandler{
			cmd:  cmd,
			args: args,
			conf: config,
		},
	}

	if err := generator.init(); err != nil {
		return nil, err
	}

	return generator, nil
}

type Generator struct {
	BaseHandler

	combinationDefinitionManager *service.CombinationStarDefinitionManager
	planDefinitionManager        *service.PlanDefinitionManager
	planExporter                 *plan.Exporter
}

func (g *Generator) Combination() {
	combinationDefinitionName := g.getArg(0, "combination")
	if combinationDefinitionName == "" {
		g.cmd.PrintErr(msgCombinationDefinitionNameMissing)
		return
	}

	planDefinitionName := g.getArg(1, "plan")
	if planDefinitionName == "" {
		g.cmd.PrintErr(msgPlanDefinitionNameMissing)
		return
	}

	combinationDefinitionScript, err := g.combinationDefinitionManager.GetScript(combinationDefinitionName)
	if err != nil {
		g.cmd.Println("Error getting combination definition:", err)
		return
	}
	planDefinitionScript, err := g.planDefinitionManager.GetFile(planDefinitionName)
	if err != nil {
		g.cmd.Println("Error getting plan definition:", err)
		return
	}

	// measure execution time
	start := time.Now()
	newPlan, err := plan.NewBuilderFromStarConfig(combinationDefinitionScript, planDefinitionScript).Build()
	if err != nil {
		g.cmd.Println("Error generating combination:", err)
		return
	}
	g.cmd.Println("Plan generated with", combinationDefinitionName, "and", planDefinitionName, "in", time.Since(start))

	start = time.Now()
	if err := g.planExporter.Export(newPlan); err != nil {
		g.cmd.Println("Error exporting plan:", err)
		return
	}
	newPlan = nil
	g.cmd.Println("Plan exported in", time.Since(start))

}

func (g *Generator) init() error {
	err := g.initFolders()
	if err != nil {
		return err
	}

	g.combinationDefinitionManager = service.NewCombinationStarDefinitionManager(g.conf.DefinitionFolder())
	g.planDefinitionManager = service.NewPlanDefinitionManager(g.conf.PlanFolder())
	g.planExporter = plan.NewExporter(g.conf.CombinationFolder(), g.conf.StorageFolder())

	return nil
}

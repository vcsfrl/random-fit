package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/plan"
	"github.com/vcsfrl/random-fit/internal/service"
	"net/http"
	"time"
)

import _ "github.com/mkevac/debugcharts"

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

	g.cmd.Println("Generating plan with combination definition ", combinationDefinitionName, "and plan definition", planDefinitionName)
	g.startMonitor()

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
	newPlan := plan.NewBuilderFromStarConfig(combinationDefinitionScript, planDefinitionScript).Generate()

	start := time.Now()
	if err := g.planExporter.ExportGenerator(newPlan); err != nil {
		g.cmd.Println("Error exporting plan:", err)
		return
	}
	g.cmd.Println("UserPlan generated and exported in", time.Since(start))
}

func (g *Generator) startMonitor() {
	g.cmd.Println("Starting debug chart server on port", g.conf.DebugChartPort)
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%s", g.conf.DebugChartPort)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			g.cmd.Println("Error debug chart server:", err)
		}
	}()

	go func() {
		<-g.cmd.Context().Done()
		if err := server.Shutdown(context.Background()); err != nil {
			g.cmd.Println("Error shutting down server:", err)
			return
		}

		g.cmd.Println("Debug chart server shut down gracefully.")
	}()
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

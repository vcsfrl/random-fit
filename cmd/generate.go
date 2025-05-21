package cmd

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/plan"
	"github.com/vcsfrl/random-fit/internal/service"
	"net/http"
	"sync"
	"time"
)

import _ "github.com/mkevac/debugcharts"

const defaultWorkers = 2

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
	logger                       zerolog.Logger
	workers                      int
}

func (g *Generator) Combination() {
	combinationDefinitionName := g.getArg(0, "combination")
	if combinationDefinitionName == "" {
		g.cmd.PrintErrln(msgCombinationDefinitionNameMissing)
		return
	}

	planDefinitionName := g.getArg(1, "plan")
	if planDefinitionName == "" {
		g.cmd.PrintErrln(msgPlanDefinitionNameMissing)
		return
	}

	g.cmd.Println("Generating plan with combination definition", combinationDefinitionName, "and plan definition", planDefinitionName)

	combinationDefinitionScript, err := g.combinationDefinitionManager.GetScript(combinationDefinitionName)
	if err != nil {
		g.cmd.PrintErrln("Error getting combination definition:", err)
		return
	}
	planDefinitionScript, err := g.planDefinitionManager.GetFile(planDefinitionName)
	if err != nil {
		g.cmd.PrintErrln("Error getting plan definition:", err)
		return
	}

	g.startMonitor()

	planGenerator := plan.NewBuilderFromStarConfig(combinationDefinitionScript, planDefinitionScript).Generate()

	start := time.Now()
	if g.export(planGenerator) {
		return
	}
	g.logger.Info().Dur("duration", time.Since(start)).Msg("UserPlan generated and exported")
	g.cmd.Println("UserPlan generated and exported in", time.Since(start))
}

func (g *Generator) export(planGenerator chan *plan.PlannedCombination) bool {
	wg := sync.WaitGroup{}

	g.logger.Info().Msgf("Starting %d workers to export plans", g.nrWorkers())
	for i := 0; i < g.nrWorkers(); i++ {
		wg.Add(1)
		g.logger.Info().Msgf("Starting worker %d", i)
		go func(i int) {
			defer func() {
				g.logger.Info().Msgf("Finished worker %d", i)
				wg.Done()
			}()

			if err := g.planExporter.ExportGenerator(planGenerator); err != nil {
				g.logger.Error().Err(err).Msgf("Error exporting plan in worker %d", i)
				g.cmd.Println("Error exporting plan:", err)
			}

			g.logger.Info().Msgf("Worker %d finished exporting plan", i)
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()
	return false
}

func (g *Generator) startMonitor() {
	g.logger.Info().Msgf("Starting debug chart server on port %s", g.conf.DebugChartPort)
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%s", g.conf.DebugChartPort)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			g.logger.Error().Err(err).Msg("Server error.")
		}
	}()

	go func() {
		<-g.cmd.Context().Done()
		if err := server.Shutdown(context.Background()); err != nil {
			g.logger.Error().Err(err).Msg("Error shutting down debug chart server.")
			return
		}

		g.logger.Info().Msg("Debug chart server shut down gracefully.")
	}()
}

func (g *Generator) nrWorkers() int {
	if g.workers == 0 {
		return defaultWorkers
	}

	return g.workers
}

func (g *Generator) init() error {
	err := g.initFolders()
	if err != nil {
		return err
	}
	g.logger = NewLogger()

	g.combinationDefinitionManager = service.NewCombinationStarDefinitionManager(g.conf.DefinitionFolder())
	g.planDefinitionManager = service.NewPlanDefinitionManager(g.conf.PlanFolder())
	g.planExporter = plan.NewExporter(g.conf.CombinationFolder(), g.conf.StorageFolder())

	return nil
}

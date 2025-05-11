package shell

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/vcsfrl/random-fit/internal/plan"
	"time"
)

func (s *Shell) generateCombination() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     "generate-combination",
		Help:     "Generate a combination",
		LongHelp: "Generate a combination.\nUsage: <shell> generate-combination <combination_definition_name> <plan_definition_name>",
		Func: func(c *ishell.Context) {
			c.Println("Generating combination...\n")

			if len(c.Args) < 2 {
				c.Println(messagePrompt + "Error: combination definition name and plan definition name are required.")
				return
			}

			combinationDefinitionName := c.Args[0]
			planDefinitionName := c.Args[1]
			combinationDefinition, err := s.getCombinationDefinitionManager().GetScript(combinationDefinitionName)
			if err != nil {
				c.Println(messagePrompt+"Error getting combination definition:", err)
				return
			}
			planDefinition, err := s.getPlanDefinitionManager().GetFile(planDefinitionName)
			if err != nil {
				c.Println(messagePrompt+"Error getting plan definition:", err)
				return
			}

			// measure execution time
			start := time.Now()
			newPlan, err := plan.NewBuilderFromStarConfig(combinationDefinition, planDefinition).Build()
			if err != nil {
				c.Println(messagePrompt+"Error generating combination:", err)
				return
			}
			c.Println(messagePrompt+"Plan generated with", combinationDefinitionName, "and", planDefinitionName, "in", time.Since(start), "\n")

			start = time.Now()
			if err := s.getExporter().Export(newPlan); err != nil {
				c.Println(messagePrompt+"Error exporting plan:", err)
				return
			}
			newPlan = nil
			c.Println(messagePrompt+"Plan exported in", time.Since(start), "\n")
		},
	}
}

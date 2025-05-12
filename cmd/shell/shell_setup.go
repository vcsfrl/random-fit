package shell

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/vcsfrl/random-fit/internal/plan"
	"os"
	"time"
)

const msgExiting = "Exiting..."

func (s *Shell) exitCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "exit",
		Help: "exit the program",
		Func: func(c *ishell.Context) {
			c.Println(msgPrompt + msgExiting)
			_ = s.Close()
			time.Sleep(100 * time.Millisecond)
		},
	}
}

func (s *Shell) execCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     "exec",
		Help:     "Execute a command non-interactively",
		LongHelp: "Execute a command non-interactively.\nUsage: <shell> exec <command>",
	}
}

func (s *Shell) interruptFunc(c *ishell.Context, count int, line string) {
	if count >= 2 {
		s.ctxCancel()
		s.shell.Close()

		time.Sleep(100 * time.Millisecond)
		c.Println("Interrupted")

		os.Exit(1)
	}
	c.Println("Input Ctrl-c once more to exit")
}

func (s *Shell) getCombinationDefinitionManager() *CombinationStarDefinitionManager {
	if s.combinationDefinitionManager == nil {
		s.combinationDefinitionManager = NewCombinationStarDefinitionManager(s.definitionFolder)
	}

	return s.combinationDefinitionManager
}

func (s *Shell) getPlanDefinitionManager() *PlanDefinitionManager {
	if s.planDefinitionManager == nil {
		s.planDefinitionManager = NewPlanDefinitionManager(s.planFolder)
	}

	return s.planDefinitionManager
}

func (s *Shell) getExporter() *plan.Exporter {
	if s.exporter == nil {
		s.exporter = plan.NewExporter(s.combinationFolder, s.storageFolder)
	}

	return s.exporter
}

func (s *Shell) createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}

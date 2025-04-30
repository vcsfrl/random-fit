package shell

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"io"
	"os"
	"os/exec"
)

const prompt = ">>> "
const messagePrompt = "-> "

type Shell struct {
	shell  *ishell.Shell
	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer

	combinationDefinitionManager *CombinationStarDefinitionManager
	planDefinitionManager        *PlanDefinitionManager
}

func New() *Shell {
	newShell := &Shell{}
	newShell.stdin = os.Stdin
	newShell.stdout = os.Stdout
	newShell.stderr = os.Stderr

	datatFolder := os.Getenv("RF_DATA_FOLDER")
	newShell.combinationDefinitionManager = NewCombinationStarDefinitionManager(datatFolder + "/definition")
	newShell.planDefinitionManager = NewPlanDefinitionManager(datatFolder + "/plan")

	newShell.init()

	return newShell
}

func (s *Shell) Run() {
	s.shell.Println("\n==================")
	s.shell.Println("=== Random-fit ===")
	s.shell.Println("==================\n")

	if len(os.Args) > 1 && os.Args[1] == "exec" {
		err := s.shell.Process(os.Args[2:]...)
		if err != nil {
			s.shell.Println("Error:", err)
		}
	} else {
		s.shell.Run()
	}
}

func (s *Shell) init() {
	s.shell = ishell.NewWithConfig(&readline.Config{
		Prompt: prompt,
		Stdin:  s.stdin,
		Stdout: s.stdout,
		Stderr: s.stderr,
	})

	s.shell.AddCmd(&ishell.Cmd{
		Name:     "exec",
		Help:     "Execute a command non-interactively",
		LongHelp: "Execute a command non-interactively.\nUsage: <shell> exec <command>",
	})

	s.shell.AddCmd(s.combinationDefinitionCmd())
	s.shell.AddCmd(s.planDefinitionCmd())
	s.shell.AddCmd(s.generateCode())
}

func (s *Shell) editScript(scriptName string, filetype string) error {
	cmd := exec.Command(os.Getenv("EDITOR"), "-filetype", filetype, scriptName)
	cmd.Stdin = s.stdin
	cmd.Stdout = s.stdout
	cmd.Stderr = s.stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

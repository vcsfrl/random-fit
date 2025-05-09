package shell

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"github.com/vcsfrl/random-fit/internal/plan"
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
	exporter                     *plan.Exporter

	definitionFolder  string
	planFolder        string
	storageFolder     string
	combinationFolder string
}

func New() *Shell {
	newShell := &Shell{}
	newShell.stdin = os.Stdin
	newShell.stdout = os.Stdout
	newShell.stderr = os.Stderr

	datatFolder := os.Getenv("RF_DATA_FOLDER")
	if datatFolder != "" {
		if err := createFolder(datatFolder); err != nil {
			newShell.shell.Println(messagePrompt+"Error creating data folder:", err)
		}

		newShell.definitionFolder = datatFolder + "/definition"
		newShell.combinationFolder = datatFolder + "/combination"
		newShell.planFolder = datatFolder + "/plan"
		newShell.storageFolder = datatFolder + "/storage"

		// create folders if they do not exist
		if err := createFolder(newShell.definitionFolder); err != nil {
			newShell.shell.Println(messagePrompt+"Error creating definition folder:", err)
		}
		if err := createFolder(newShell.combinationFolder); err != nil {
			newShell.shell.Println(messagePrompt+"Error creating combination folder:", err)
		}
		if err := createFolder(newShell.planFolder); err != nil {
			newShell.shell.Println(messagePrompt+"Error creating plan folder:", err)
		}
		if err := createFolder(newShell.storageFolder); err != nil {
			newShell.shell.Println(messagePrompt+"Error creating storage folder:", err)
		}
	}

	newShell.init()

	return newShell
}

func (s *Shell) Run() {
	s.shell.Println("\n==================")
	s.shell.Println("=== Random-fit ===")
	s.shell.Println("==================\n")

	defer func() {
		// recover from panic if one occurred.
		if err := recover(); err != nil {
			s.shell.Println(messagePrompt+"Error:", err)

		}
		s.shell.Stop()
	}()

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
	s.shell.AddCmd(s.generateCombination())
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

func createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}

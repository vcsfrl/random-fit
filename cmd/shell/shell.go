package shell

import (
	"context"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"github.com/vcsfrl/random-fit/internal/plan"
	"io"
	"os"
)

const prompt = ">>> "
const messagePrompt = "-> "
const separator = "=================="
const welcomeMessage = "=== Random-fit ==="

type Shell struct {
	shell       *ishell.Shell
	stdin       io.ReadCloser
	stdout      io.Writer
	stdinWriter io.Writer

	stderr                       io.Writer
	combinationDefinitionManager *CombinationStarDefinitionManager
	planDefinitionManager        *PlanDefinitionManager

	exporter         *plan.Exporter
	definitionFolder string
	planFolder       string
	storageFolder    string

	combinationFolder string
	ctx               context.Context
	ctxCancel         context.CancelFunc
}

func New() *Shell {
	newShell := BuildNew()
	newShell.Init()

	return newShell
}

func BuildNew() *Shell {
	newShell := &Shell{}
	newShell.stdin = os.Stdin
	newShell.stdout = os.Stdout
	newShell.stderr = os.Stderr
	newShell.stdinWriter = os.Stdin
	newShell.ctx, newShell.ctxCancel = context.WithCancel(context.Background())

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
	return newShell
}

func (s *Shell) Init() {
	s.shell = ishell.NewWithConfig(&readline.Config{
		Prompt:      prompt,
		Stdin:       s.stdin,
		StdinWriter: s.stdinWriter,
		Stdout:      s.stdout,
		Stderr:      s.stderr,
	})

	s.shell.DeleteCmd("exit")
	s.shell.AddCmd(s.exitCmd())
	s.shell.Interrupt(s.interruptFunc)

	s.shell.AddCmd(s.execCmd())
	s.shell.AddCmd(s.combinationDefinitionCmd())
	s.shell.AddCmd(s.planDefinitionCmd())
	s.shell.AddCmd(s.generateCode())
	s.shell.AddCmd(s.generateCombination())
}

func (s *Shell) Run() {
	s.shell.Println(separator)
	s.shell.Println(welcomeMessage)
	s.shell.Println(separator, "\n")

	s.runTrace()

	defer func() {
		// handle panic.
		if err := recover(); err != nil {
			s.shell.Println(messagePrompt+"Error:", err)
		}
		s.shell.Close()
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

func (s *Shell) Close() error {
	s.ctxCancel()
	s.shell.Close()

	return nil
}

//
//func (s *Shell) RunCommand(command string) {
//	if _, err := s.stdinWriter.Write([]byte(command + "\n")); err != nil {
//		s.shell.Println(messagePrompt+"Error writing command to stdin:", err)
//	}
//}

func createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}

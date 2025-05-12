package shell

import (
	"context"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"github.com/vcsfrl/random-fit/internal/plan"
	"io"
	"os"
)

const shellPrompt = ">>> "
const msgPrompt = "-> "
const msgSeparator = "=================="
const msgWelcomeMessage = "=== Random-fit ==="

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
	newShell := &Shell{}
	newShell.stdin = os.Stdin
	newShell.stdout = os.Stdout
	newShell.stderr = os.Stderr
	newShell.stdinWriter = os.Stdin

	datatFolder := os.Getenv("RF_DATA_FOLDER")
	if datatFolder != "" {
		if err := newShell.createFolder(datatFolder); err != nil {
			newShell.shell.Println(msgPrompt+"Error creating data folder:", err)
		}

		newShell.definitionFolder = datatFolder + "/definition"
		newShell.combinationFolder = datatFolder + "/combination"
		newShell.planFolder = datatFolder + "/plan"
		newShell.storageFolder = datatFolder + "/storage"

		// create folders if they do not exist
		if err := newShell.createFolder(newShell.definitionFolder); err != nil {
			newShell.shell.Println(msgPrompt+"Error creating definition folder:", err)
		}
		if err := newShell.createFolder(newShell.combinationFolder); err != nil {
			newShell.shell.Println(msgPrompt+"Error creating combination folder:", err)
		}
		if err := newShell.createFolder(newShell.planFolder); err != nil {
			newShell.shell.Println(msgPrompt+"Error creating plan folder:", err)
		}
		if err := newShell.createFolder(newShell.storageFolder); err != nil {
			newShell.shell.Println(msgPrompt+"Error creating storage folder:", err)
		}
	}

	return newShell
}

func (s *Shell) Run() {
	s.init()
	s.shell.Println(msgSeparator)
	s.shell.Println(msgWelcomeMessage)
	s.shell.Println(msgSeparator, "\n")

	s.runTrace()

	defer func() {
		// handle panic.
		if err := recover(); err != nil {
			s.shell.Println(msgPrompt+"Error:", err)
		}
		_ = s.Close()
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

func (s *Shell) SendCommand(command string) {
	if _, err := io.WriteString(s.stdinWriter, command+"\n"); err != nil {
		s.shell.Println(msgPrompt+"Error writing command to stdin:", err)
	}
}

func (s *Shell) Close() error {
	if s.ctxCancel != nil {
		s.ctxCancel()
	}

	if s.shell != nil {
		s.shell.Close()
		return nil
	}

	return nil
}

func (s *Shell) init() {
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	s.shell = ishell.NewWithConfig(&readline.Config{
		Prompt:      shellPrompt,
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

package shell

import (
	"github.com/abiosoft/ishell/v2"
	"os"
)

type Shell struct {
	shell *ishell.Shell
}

func New() *Shell {
	newShell := &Shell{}
	newShell.init()

	return newShell
}

func (s *Shell) Run() {
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
	s.shell = ishell.New()
	s.shell.Println("Random-fit")
	s.shell.AddCmd(s.definitionCmd())
	s.shell.AddCmd(&ishell.Cmd{
		Name:     "exec",
		Help:     "Execute a command non-interactively",
		LongHelp: "Execute a command non-interactively.\nUsage: <shell> exec <command>",
	})
}

func (s *Shell) definitionCmd() *ishell.Cmd {
	listDefinition := &ishell.Cmd{
		Name: "list",
		Help: "List definitions",
		Func: func(c *ishell.Context) {
			c.Println("List definitions")
		},
	}

	definition := &ishell.Cmd{
		Name: "definition",
		Help: "Manage definitions",
		Func: func(c *ishell.Context) {
			listDefinition.Func(c)
		},
	}

	definition.AddCmd(listDefinition)

	return definition
}

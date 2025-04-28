package shell

import (
	"bytes"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	textTemplate "text/template"
)

const definitionSkeleton = `package shell

// definitionTemplate is a template for a definition file
// This is a generated file. Do not edit.
var definitionTemplate = {{.}}`

type Shell struct {
	shell  *ishell.Shell
	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer

	definitionManager *StarDefinitionManager
}

func New() *Shell {
	newShell := &Shell{}
	newShell.init()
	newShell.stdin = os.Stdin
	newShell.stdout = os.Stdout
	newShell.stderr = os.Stderr

	datatFolder := os.Getenv("RF_DATA_FOLDER")
	newShell.definitionManager = NewStarDefinitionManager(datatFolder + "/definition")

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
	s.shell = ishell.NewWithConfig(&readline.Config{
		Prompt: ">>> ",
		Stdin:  s.stdin,
		Stdout: s.stdout,
		Stderr: s.stderr,
	})
	s.shell.Println("==============")
	s.shell.Println("= Random-fit =")
	s.shell.Println("==============\n")
	s.shell.AddCmd(s.definitionCmd())
	s.shell.AddCmd(&ishell.Cmd{
		Name:     "exec",
		Help:     "Execute a command non-interactively",
		LongHelp: "Execute a command non-interactively.\nUsage: <shell> exec <command>",
	})
	s.shell.AddCmd(s.generateCode())
}

func (s *Shell) definitionCmd() *ishell.Cmd {
	listDefinition := &ishell.Cmd{
		Name: "list",
		Help: "List definitions",
		Func: func(c *ishell.Context) {
			c.Println("Definitions:")

			definitions, err := s.definitionManager.List()
			if err != nil {
				c.Println("-> Error listing definition:", err)
				return
			}

			if len(definitions) == 0 {
				c.Println("-> No definitions found.")
				return
			}

			for _, definition := range definitions {
				c.Println(" - ", definition)
			}
		},
	}

	newDefinition := &ishell.Cmd{
		Name:     "new",
		Help:     "Create a new definition",
		LongHelp: "Create a new definition.\nUsage: <shell> definition new <definition_name>",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("-> Error: definition name is required.")
				return
			}

			err := s.definitionManager.New(c.Args[0])
			if err != nil {
				c.Println("-> Error new definition:", err)
				return
			}

			c.Println("-> Definition created:", c.Args[0], "\n")
		},
	}

	editDefinition := &ishell.Cmd{
		Name:     "edit",
		Help:     "Edit definition",
		LongHelp: "Edit a definition.",
		Func: func(c *ishell.Context) {
			definitions, err := s.definitionManager.List()
			if err != nil {
				c.Println("-> Error getting definitions list:", err)
				return
			}

			choice := c.MultiChoice(definitions, "Select a definition to edit:")

			scriptName, err := s.definitionManager.GetScript(definitions[choice])
			if err != nil {
				c.Println("-> Error getting definition script:", err)
				return
			}

			cmd := exec.Command(os.Getenv("EDITOR"), scriptName)
			cmd.Stdin = s.stdin
			cmd.Stdout = s.stdout
			cmd.Stderr = s.stderr

			if err = cmd.Start(); err != nil {
				c.Println("-> Error starting editor:", err)
				return
			}
			if err := cmd.Wait(); err != nil {
				c.Println("-> Error waiting for editor:", err)
				return
			}

			c.Println("-> Definition edited:", definitions[choice], "\n")
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
	definition.AddCmd(newDefinition)
	definition.AddCmd(editDefinition)

	return definition
}

func (s *Shell) generateCode() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "generate-code",
		Help: "Generate code",
		Func: func(c *ishell.Context) {
			c.Println("Generating helper code...\n")

			baseFolder := os.Getenv("RF_BASE_FOLDER")

			t := textTemplate.Must(textTemplate.New("template.render_text").Parse(definitionSkeleton))

			//create a file in shell/ folder
			fileName := filepath.Join(baseFolder, "cmd", "shell", "definition_template.go")
			// remove the file if it exists
			if err := os.Remove(fileName); err != nil && !os.IsNotExist(err) {
				c.Println("Error:", err)
				return
			}

			// get content of star definition template
			content, err := os.ReadFile(filepath.Join(baseFolder, "internal", "combination", "template", "script.star"))
			if err != nil {
				c.Println("Error:", err)
				return
			}

			buff := &bytes.Buffer{}
			if err := t.Execute(buff, "`"+string(content)+"`"); err != nil {
				c.Println("Error:", err)
				return
			}

			if err := os.WriteFile(fileName, buff.Bytes(), 0644); err != nil {
				c.Println("Error:", err)
			}

			c.Println("Code generated in", fileName, "\n")
		},
	}
}

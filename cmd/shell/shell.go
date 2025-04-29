package shell

import (
	"bytes"
	"encoding/json"
	"github.com/abiosoft/ishell/v2"
	"github.com/abiosoft/readline"
	"github.com/charmbracelet/glamour"
	"github.com/vcsfrl/random-fit/internal/combination"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	textTemplate "text/template"
)

const definitionSkeleton = `package shell

// This is a generated file. Do not edit!

// definitionTemplate is a template for a definition file
var definitionTemplate = {{.}}`

const prompt = ">>> "
const messagePrompt = "-> "

type Shell struct {
	shell  *ishell.Shell
	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer

	definitionManager *StarDefinitionManager
}

func New() *Shell {
	newShell := &Shell{}
	newShell.stdin = os.Stdin
	newShell.stdout = os.Stdout
	newShell.stderr = os.Stderr

	datatFolder := os.Getenv("RF_DATA_FOLDER")
	newShell.definitionManager = NewStarDefinitionManager(datatFolder + "/definition")

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
				c.Println(messagePrompt+"Error listing definition:", err)
				return
			}

			if len(definitions) == 0 {
				c.Println(messagePrompt + "No definitions found.")
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
				c.Println(messagePrompt + "Error: definition name is required.")
				return
			}

			err := s.definitionManager.New(c.Args[0])
			if err != nil {
				c.Println(messagePrompt+"Error new definition:", err)
				return
			}
			c.Println(messagePrompt+"Definition created:", c.Args[0], "\n")

			if err := s.editDefinition(c.Args[0]); err != nil {
				c.Println(messagePrompt+"Error editing definition:", err)
				return
			}
		},
	}

	editDefinition := &ishell.Cmd{
		Name:     "edit",
		Help:     "Edit definition",
		LongHelp: "Edit a definition.",
		Func: func(c *ishell.Context) {
			definitions, err := s.definitionManager.List()
			if err != nil {
				c.Println(messagePrompt+"Error getting definitions list:", err)
				return
			}
			choice := c.MultiChoice(definitions, "Select a definition to edit:")

			if err := s.editDefinition(definitions[choice]); err != nil {
				c.Println(messagePrompt+"Error editing definition:", err)
				return
			}

			c.Println(messagePrompt+"Definition edited:", definitions[choice], "\n")
		},
	}

	viewDefinition := &ishell.Cmd{
		Name:     "view",
		Help:     "View definition",
		LongHelp: "View a definition.",
		Func: func(c *ishell.Context) {
			_ = c.ClearScreen()
			definitions, err := s.definitionManager.List()
			if err != nil {
				c.Println(messagePrompt+"Error getting definitions list:", err)
				return
			}
			choice := c.MultiChoice(definitions, "Select a definition to view:")

			viewCombination, err := s.definitionManager.Build(definitions[choice])
			if err != nil {
				c.Println(messagePrompt+"Error building definition:", err)
				return
			}

			c.Println(messagePrompt+"Combination:", viewCombination.Details)
			c.Println(messagePrompt+"Definition ID:", viewCombination.DefinitionID)

			for dataType, data := range viewCombination.Data {
				c.Println(messagePrompt+"Definition view:", dataType)
				c.Println(messagePrompt + "====================================")
				err := s.printCombination(c, data)
				if err != nil {
					c.Println(messagePrompt+"Error viewing data:", err)
					return
				}
				c.Println(messagePrompt + "====================================")
			}

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
	definition.AddCmd(viewDefinition)

	return definition
}

func (s *Shell) editDefinition(definition string) error {
	scriptName, err := s.definitionManager.GetScript(definition)
	if err != nil {
		return err
	}

	if err := s.editDefinitionScript(scriptName); err != nil {
		return err
	}

	return nil
}

func (s *Shell) editDefinitionScript(scriptName string) error {
	cmd := exec.Command(os.Getenv("EDITOR"), scriptName)
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

func (s *Shell) printCombination(c *ishell.Context, data *combination.Data) error {
	switch data.Type {
	case combination.DataTypeJson:
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, data.Data.Bytes(), "", "  ")
		if err != nil {
			return err
		}
		c.Println(prettyJSON.String())
	case combination.DataTypeMd:
		out, err := glamour.Render(data.Data.String(), "dark")
		if err != nil {
			return err
		}
		c.Println(out)
	default:
		c.Println(data.Data.String())
	}

	return nil
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
				c.Println(messagePrompt+"Error:", err)
				return
			}

			// get content of star definition template
			content, err := os.ReadFile(filepath.Join(baseFolder, "internal", "combination", "template", "script.star"))
			if err != nil {
				c.Println(messagePrompt+"Error:", err)
				return
			}

			buff := &bytes.Buffer{}
			if err := t.Execute(buff, "`"+string(content)+"`"); err != nil {
				c.Println(messagePrompt+"Error:", err)
				return
			}

			if err := os.WriteFile(fileName, buff.Bytes(), 0644); err != nil {
				c.Println(messagePrompt+"Error:", err)
			}

			c.Println(messagePrompt+"Code generated in", fileName, "\n")
		},
	}
}

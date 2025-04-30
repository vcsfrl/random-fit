package shell

import (
	"bytes"
	"encoding/json"
	"github.com/abiosoft/ishell/v2"
	"github.com/charmbracelet/glamour"
	"github.com/vcsfrl/random-fit/internal/combination"
)

func (s *Shell) combinationDefinitionCmd() *ishell.Cmd {
	listDefinition := &ishell.Cmd{
		Name: "list",
		Help: "List definitions",
		Func: func(c *ishell.Context) {
			c.Println("Definitions:")

			definitions, err := s.combinationDefinitionManager.List()
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

			err := s.combinationDefinitionManager.New(c.Args[0])
			if err != nil {
				c.Println(messagePrompt+"Error new definition:", err)
				return
			}
			c.Println(messagePrompt+"Definition created:", c.Args[0], "\n")

			if err := s.editCombinationDefinition(c.Args[0]); err != nil {
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
			definitions, err := s.combinationDefinitionManager.List()
			if err != nil {
				c.Println(messagePrompt+"Error getting definitions list:", err)
				return
			}
			choice := c.MultiChoice(definitions, "Select a definition to edit:")

			if err := s.editCombinationDefinition(definitions[choice]); err != nil {
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
			definitions, err := s.combinationDefinitionManager.List()
			if err != nil {
				c.Println(messagePrompt+"Error getting definitions list:", err)
				return
			}
			choice := c.MultiChoice(definitions, "Select a definition to view:")

			viewCombination, err := s.combinationDefinitionManager.Build(definitions[choice])
			if err != nil {
				c.Println(messagePrompt+"Error building definition:", err)
				return
			}

			c.Println(messagePrompt+"Combination:", viewCombination.Details)
			c.Println(messagePrompt+"Definition ID:", viewCombination.DefinitionID)

			for dataType, data := range viewCombination.Data {
				c.Println(messagePrompt+"Definition view:", dataType)
				c.Println(messagePrompt + "====================================")
				err := s.printCombinationDefinition(c, data)
				if err != nil {
					c.Println(messagePrompt+"Error viewing data:", err)
					return
				}
				c.Println(messagePrompt + "====================================")
			}

		},
	}

	definition := &ishell.Cmd{
		Name: "combination-definition",
		Help: "Manage combination definitions",
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

func (s *Shell) editCombinationDefinition(definition string) error {
	scriptName, err := s.combinationDefinitionManager.GetScript(definition)
	if err != nil {
		return err
	}

	if err := s.editScript(scriptName, "python"); err != nil {
		return err
	}

	return nil
}

func (s *Shell) printCombinationDefinition(c *ishell.Context, data *combination.Data) error {
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

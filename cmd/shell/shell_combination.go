package shell

import (
	"bytes"
	"encoding/json"
	"github.com/abiosoft/ishell/v2"
	"github.com/charmbracelet/glamour"
	"github.com/vcsfrl/random-fit/internal/combination"
)

const cmdCombinationDefinitionName = "combination-definition"
const cmdCombinationDefinitionHelp = "Manage combination definitions"
const msgNoDefinitions = "No definitions found."

func (s *Shell) combinationDefinitionCmd() *ishell.Cmd {
	listDefinition := &ishell.Cmd{
		Name: "list",
		Help: "List definitions",
		Func: func(c *ishell.Context) {
			c.Println("Definitions:")
			definitions, err := s.getCombinationDefinitionManager().List()
			if err != nil {
				c.Println(msgPrompt+"Error listing definition:", err)
				return
			}

			if len(definitions) == 0 {
				c.Println(msgPrompt + msgNoDefinitions)
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
				c.Println(msgPrompt + "Error: definition name is required.")
				return
			}

			err := s.getCombinationDefinitionManager().New(c.Args[0])
			if err != nil {
				c.Println(msgPrompt+"Error new definition:", err)
				return
			}
			c.Println(msgPrompt+"Definition created:", c.Args[0], "\n")

			if err := s.editCombinationDefinition(c.Args[0]); err != nil {
				c.Println(msgPrompt+"Error editing definition:", err)
				return
			}
		},
	}

	editDefinition := &ishell.Cmd{
		Name:     "edit",
		Help:     "Edit definition",
		LongHelp: "Edit a definition.",
		Func: func(c *ishell.Context) {
			selectedDefinition, err := s.getSelectedDefinition(c)
			if err != nil {
				c.Println(msgPrompt+"Error getting definition:", err)
				return
			}

			if err := s.editCombinationDefinition(selectedDefinition); err != nil {
				c.Println(msgPrompt+"Error editing definition:", err)
				return
			}

			c.Println(msgPrompt+"Definition edited:", selectedDefinition, "\n")
		},
	}

	viewDefinition := &ishell.Cmd{
		Name:     "view",
		Help:     "View definition",
		LongHelp: "View a definition.",
		Func: func(c *ishell.Context) {
			_ = c.ClearScreen()
			selectedDefinition, err := s.getSelectedDefinition(c)
			if err != nil {
				c.Println(msgPrompt+"Error getting definition:", err)
				return
			}

			viewCombination, err := s.getCombinationDefinitionManager().Build(selectedDefinition)
			if err != nil {
				c.Println(msgPrompt+"Error building definition:", err)
				return
			}

			c.Println(msgPrompt+"Combination:", viewCombination.Details)
			c.Println(msgPrompt+"Definition ID:", viewCombination.DefinitionID)

			for dataType, data := range viewCombination.Data {
				c.Println(msgPrompt+"Definition view:", dataType)
				c.Println(msgPrompt + msgSeparator + msgSeparator)
				err := s.printCombinationDefinition(c, data)
				if err != nil {
					c.Println(msgPrompt+"Error viewing data:", err)
					return
				}
				c.Println(msgPrompt + msgSeparator + msgSeparator)
			}
		},
	}

	definition := &ishell.Cmd{
		Name: cmdCombinationDefinitionName,
		Help: cmdCombinationDefinitionHelp,
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

func (s *Shell) getSelectedDefinition(c *ishell.Context) (string, error) {
	var selectedDefinition string
	if len(c.Args) == 0 {
		definitions, err := s.getCombinationDefinitionManager().List()
		if err != nil {
			return "", err
		}
		choice := c.MultiChoice(definitions, "Select a definition to edit:")

		selectedDefinition = definitions[choice]
	} else {
		selectedDefinition = c.Args[0]
		if _, err := s.getCombinationDefinitionManager().GetScript(selectedDefinition); err != nil {
			return "", err
		}
	}

	return selectedDefinition, nil
}

func (s *Shell) editCombinationDefinition(definition string) error {
	scriptName, err := s.getCombinationDefinitionManager().GetScript(definition)
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

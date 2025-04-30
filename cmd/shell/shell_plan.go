package shell

import "github.com/abiosoft/ishell/v2"

func (s *Shell) planDefinitionCmd() *ishell.Cmd {
	listDefinition := &ishell.Cmd{
		Name: "list",
		Help: "List plan definitions",
		Func: func(c *ishell.Context) {
			c.Println("Plan definitions:")

			definitions, err := s.getPlanDefinitionManager().List()
			if err != nil {
				c.Println(messagePrompt+"Error listing plan definitions:", err)
				return
			}

			if len(definitions) == 0 {
				c.Println(messagePrompt + "No plan definitions found.")
				return
			}

			for _, definition := range definitions {
				c.Println(" - ", definition)
			}

		},
	}

	newDefinition := &ishell.Cmd{
		Name:     "new",
		Help:     "Create a new plan definition",
		LongHelp: "Create a new plan definition.\nUsage: <shell> plan-definition new <definition_name>",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println(messagePrompt + "Error: plan definition name is required.")
				return
			}

			err := s.getPlanDefinitionManager().New(c.Args[0])
			if err != nil {
				c.Println(messagePrompt+"Error creating new plan definition:", err)
				return
			}
			c.Println(messagePrompt+"Plan definition created:", c.Args[0], "\n")

			if err := s.editPlanDefinition(c.Args[0]); err != nil {
				c.Println(messagePrompt+"Error editing plan definition:", err)
				return
			}
		},
	}

	editDefinition := &ishell.Cmd{
		Name:     "edit",
		Help:     "Edit plan definition",
		LongHelp: "Edit a plan definition.\nUsage: <shell> plan-definition edit",
		Func: func(c *ishell.Context) {
			definitions, err := s.getPlanDefinitionManager().List()
			if err != nil {
				c.Println(messagePrompt+"Error getting plan definitions list:", err)
				return
			}
			choice := c.MultiChoice(definitions, "Select a definition to edit:")

			if err := s.editPlanDefinition(definitions[choice]); err != nil {
				c.Println(messagePrompt+"Error editing plan definition:", err)
				return
			}
		},
	}

	definition := &ishell.Cmd{
		Name: "plan-definition",
		Help: "Manage combination definitions",
		Func: func(c *ishell.Context) {
			listDefinition.Func(c)
		},
	}

	definition.AddCmd(listDefinition)
	definition.AddCmd(newDefinition)
	definition.AddCmd(editDefinition)

	return definition
}

func (s *Shell) editPlanDefinition(definition string) error {
	scriptName, err := s.getPlanDefinitionManager().GetFile(definition)
	if err != nil {
		return err
	}

	if err := s.editScript(scriptName, "json"); err != nil {
		return err
	}

	return nil
}

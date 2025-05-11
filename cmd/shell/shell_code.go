package shell

import (
	"bytes"
	"github.com/abiosoft/ishell/v2"
	"os"
	"path/filepath"
	textTemplate "text/template"
)

const definitionSkeleton = `package shell

// This is a generated file. Do not edit!

func Init() {
	// definitionTemplate is a template for a definition file
	definitionTemplate = {{.}}
}`

func (s *Shell) generateCode() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "generate-code",
		Help: "Generate code",
		Func: func(c *ishell.Context) {
			c.Println("Generating helper code...\n")

			baseFolder := os.Getenv("RF_BASE_FOLDER")

			t := textTemplate.Must(textTemplate.New("template.render_text").Parse(definitionSkeleton))

			//create a file in shell/ folder
			fileName := filepath.Join(baseFolder, "cmd", "shell", "combination_definition_template.go")
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

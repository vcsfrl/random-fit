package internal

import (
	"bytes"
	"os"
	"path/filepath"
	textTemplate "text/template"
)

const definitionSkeleton = `package internal

// This is a generated file. Do not edit!

func init() {
	// definitionTemplate is a template for a definition file
	definitionTemplate = {{.}}
}`

func GenerateCode(c Printer, config *Config) {
	c.Println("Generating helper code...\n")

	t := textTemplate.Must(textTemplate.New("template.render_text").Parse(definitionSkeleton))

	//create a file in shell/ folder
	fileName := filepath.Join(config.BaseFolder, "cmd", "internal", "combination_definition_template.go")
	// remove the file if it exists
	if err := os.Remove(fileName); err != nil && !os.IsNotExist(err) {
		c.Println("Error:", err)
		return
	}

	// get content of star definition template
	content, err := os.ReadFile(filepath.Join(config.BaseFolder, "internal", "combination", "template", "script.star"))
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
}

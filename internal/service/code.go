package service

import (
	"bytes"
	"github.com/vcsfrl/random-fit/internal/platform/fs"
	"os"
	"path/filepath"
	textTemplate "text/template"
)

const definitionSkeleton = `package service

// This is a generated file. Do not edit!

func (dm *CombinationStarDefinitionManager) GetTemplate() string {
	// definitionTemplate is a template for a definition file
	var definitionTemplate = {{.}}

	return definitionTemplate
}`

func GenerateCode(printer Printer, config *Config) {
	printer.Println("Generating helper code...\n")

	textTmpl := textTemplate.Must(textTemplate.New("template.render_text").Parse(definitionSkeleton))

	// create a file in shell/ folder
	fileName := filepath.Join(config.BaseFolder, "internal", "service", "combination_definition_template.go")
	// remove the file if it exists
	if err := os.Remove(fileName); err != nil && !os.IsNotExist(err) {
		printer.Println("Error:", err)

		return
	}

	// get content of star definition template
	content, err := os.ReadFile(filepath.Join(config.BaseFolder, "internal", "combination", "template", "script.star"))
	if err != nil {
		printer.Println("Error:", err)

		return
	}

	buff := &bytes.Buffer{}
	if err := textTmpl.Execute(buff, "`"+string(content)+"`"); err != nil {
		printer.Println("Error:", err)

		return
	}

	if err := os.WriteFile(fileName, buff.Bytes(), fs.FilePermission); err != nil {
		printer.Println("Error:", err)
	}

	printer.Println("Code generated in", fileName, "\n")
}

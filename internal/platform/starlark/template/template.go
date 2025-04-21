package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	textTemplate "text/template"
)

var Module = &starlarkstruct.Module{
	Name: "template",
	Members: starlark.StringDict{
		"render_text": starlark.NewBuiltin("render_text", renderText),
	},
}

// renderText() is a Go function called from Starlark.
// It renders a text textTemplate with the given arguments.
func renderText(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var tpl string
	var tplJsonArgs string
	var tplGoArgs any

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "tpl", &tpl, "tplJsonArgs", &tplJsonArgs); err != nil {
		return nil, err
	}

	// Create a new textTemplate and parse the letter into it.
	t := textTemplate.Must(textTemplate.New("template.render_text").Parse(tpl))
	if err := json.Unmarshal([]byte(tplJsonArgs), &tplGoArgs); err != nil {
		return nil, fmt.Errorf("unmarshal slJson args: %w", err)
	}
	buff := &bytes.Buffer{}

	// Execute the textTemplate.
	err := t.Execute(buff, tplGoArgs)
	if err != nil {
		return nil, fmt.Errorf("execute textTemplate: %w", err)
	}

	return starlark.String(buff.String()), nil
}

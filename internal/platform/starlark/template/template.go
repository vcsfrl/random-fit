package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	textTemplate "text/template"
)

type Template struct {
	Module *starlarkstruct.Module
}

func New() *Template {
	tpl := &Template{}
	tpl.init()

	return tpl
}

func (t *Template) init() {
	t.Module = &starlarkstruct.Module{
		Name: "template",
		Members: starlark.StringDict{
			"render_text": starlark.NewBuiltin("render_text", t.renderText),
		},
	}
}

// renderText() is a Go function called from Starlark.
// It renders a text textTemplate with the given arguments.
//
//nolint:lll
func (t *Template) renderText(_ *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) { //nolint:ireturn
	var tpl, tplJSONArgs string

	var tplGoArgs any

	if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "tpl", &tpl, "tplJsonArgs", &tplJSONArgs); err != nil {
		return nil, fmt.Errorf("unpack args: %w", err)
	}

	// Create a new textTemplate and parse the letter into it.
	textTmpl := textTemplate.Must(textTemplate.New("template.render_text").Parse(tpl))

	if err := json.Unmarshal([]byte(tplJSONArgs), &tplGoArgs); err != nil {
		return nil, fmt.Errorf("unmarshal slJson args: %w", err)
	}

	buff := &bytes.Buffer{}

	// Execute the textTemplate.
	err := textTmpl.Execute(buff, tplGoArgs)
	if err != nil {
		return nil, fmt.Errorf("execute textTemplate: %w", err)
	}

	return starlark.String(buff.String()), nil
}

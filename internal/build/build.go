package build

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/model"
	"github.com/vcsfrl/random-fit/internal/tmp/platform/random"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"html/template"
	"path/filepath"
	"time"
)

var ErrBuilding = fmt.Errorf("error building combination")
var ErrBuildingScript = fmt.Errorf("%w: error starlark script", ErrBuilding)

type Builder struct {
	thread      *starlark.Thread
	builderFunc starlark.Value
	definition  *model.CombinationDefinition

	uuidFunc       func() (string, error)
	nowFunc        func() time.Time
	randomUintFunc func(min uint, max uint) (uint, error)
}

func NewBuilder(definition *model.CombinationDefinition) (*Builder, error) {
	builder := &Builder{definition: definition}
	err := builder.start()
	if err != nil {
		return nil, err
	}

	return builder, nil

}

func (bd *Builder) Build() (*model.Combination, error) {
	// Run the Starlark script from the definition to create a new combination.
	combinationData, err := starlark.Call(bd.thread, bd.builderFunc, nil, nil)
	if err != nil {
		return nil, ErrBuildingScript
	}

	// Build the template from the definition.
	base := filepath.Base(bd.definition.GoTemplate)
	templateData := template.Must(template.New(base).ParseFiles(bd.definition.GoTemplate))

	output := new(bytes.Buffer)
	if err := templateData.Execute(output, combinationData); err != nil {
		return nil, err
	}

	// Build the combination from the template and the data from the Starlark script.
	combination := &model.Combination{
		UUID:       uuid.New(),
		Definition: bd.definition,
		Data:       combinationData,
		Output:     output,
	}

	return combination, nil
}

func (bd *Builder) start() error {
	// The Thread defines the behavior of the built-in 'print' function.
	bd.thread = &starlark.Thread{
		Name:  "combination-builder",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), bd.thread, bd.definition.StarScript, nil, bd.predeclared())
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf("execution error: %w\n%s", evalErr, evalErr.Backtrace())
		}
		return fmt.Errorf("execution error: %w", err)
	}

	// Retrieve a module global.
	buildCombination, ok := globals["build_combination"]
	if !ok {
		return fmt.Errorf("missing 'build_combination' function definition in %s", bd.definition.StarScript)
	}

	bd.builderFunc = buildCombination

	return nil
}

func (bd *Builder) predeclared() starlark.StringDict {
	// uuidF() is a Go function called from Starlark.
	// It returns a new UUID.
	uuidF := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		uuidFunc := bd.getUuidFunc()
		id, err := uuidFunc()
		if err != nil {
			return nil, err

		}

		return starlark.String(id), nil
	}

	// now() is a Go function called from Starlark.
	// It returns the current time in RFC3339 format.
	now := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		now := bd.getNowFunc()()

		return starlark.String(time.Time.Format(now, time.RFC3339)), nil
	}

	// randomInt() is a Go function called from Starlark.
	// It returns multiple random values from an interval.
	randomInt := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var min uint
		var max uint
		var nr int
		var allowDuplicates = false

		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "min", &min, "max", &max, "nr", &nr, "allow_duplicates?", &allowDuplicates); err != nil {
			return nil, err
		}

		result := starlark.NewList([]starlark.Value{})
		uintFunc, err := bd.getRandomIntFunc()(min, max)
		if err != nil {
			return nil, err
		}

		for i := 0; i < nr; i++ {
			err = result.Append(starlark.MakeUint(uintFunc))
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"uuid":       starlark.NewBuiltin("uuid", uuidF),
		"now":        starlark.NewBuiltin("now", now),
		"random_int": starlark.NewBuiltin("random_int", randomInt),
	}

	return predeclared
}

func (bd *Builder) getUuidFunc() func() (string, error) {
	if bd.uuidFunc != nil {
		return bd.uuidFunc
	}

	bd.uuidFunc = func() (string, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return id.String(), nil
	}

	return bd.uuidFunc
}

func (bd *Builder) getNowFunc() func() time.Time {
	if bd.nowFunc != nil {
		return bd.nowFunc
	}

	bd.nowFunc = func() time.Time {
		return time.Now()
	}

	return bd.nowFunc
}

func (bd *Builder) getRandomIntFunc() func(min uint, max uint) (uint, error) {
	if bd.randomUintFunc != nil {
		return bd.randomUintFunc
	}

	bd.randomUintFunc = random.NewCrypto().Uint

	return bd.randomUintFunc
}

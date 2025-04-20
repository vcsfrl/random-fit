package combination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/platform/random"
	slJson "go.starlark.net/lib/json"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"slices"
	"text/template"
	"time"
)

var ErrCombinationDefinition = fmt.Errorf("error combination definition")

type StarlarkDefinition struct {
	ID         string
	Details    string
	StarScript string

	buildFunction  *starlark.Function
	thread         *starlark.Thread
	uuidFunc       func() (string, error)
	nowFunc        func() time.Time
	randomUintFunc func(min uint, max uint) (uint, error)
}

func NewCombinationDefinition(script string) (*StarlarkDefinition, error) {
	definition := &StarlarkDefinition{
		StarScript: script,
	}

	err := definition.init()
	if err != nil {
		return nil, err
	}

	return definition, nil
}

func (cd *StarlarkDefinition) CallScriptBuildFunction() (string, error) {
	combinationStarlarkData, err := starlark.Call(cd.thread, cd.buildFunction, nil, nil)
	if err != nil {
		return "", fmt.Errorf("%w: error building combination data: %w", ErrCombinationDefinition, err)
	}

	combinationDict, ok := combinationStarlarkData.(*starlark.Dict)
	if !ok {
		return "", fmt.Errorf("%w: combination data is not a dict", ErrCombinationDefinition)
	}

	return combinationDict.String(), nil
}

func (cd *StarlarkDefinition) init() error {
	// The Thread defines the behavior of the built-in 'print' function.
	cd.thread = &starlark.Thread{
		Name:  cd.Details,
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), cd.thread, cd.StarScript, nil, cd.predeclared())
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf("execution error: %w\n%s", evalErr, evalErr.Backtrace())
		}
		return fmt.Errorf("execution error: %w", err)
	}

	// Retrieve the definition from the globals.
	definition, ok := globals["definition"]
	if !ok {
		return fmt.Errorf("%w missing 'definition' Dict %s", ErrCombinationDefinition, cd.StarScript)
	}

	dictDefinition, ok := definition.(*starlark.Dict)
	if !ok {
		return fmt.Errorf("%w 'definition' must be a Dict %s", ErrCombinationDefinition, cd.StarScript)
	}

	// Retrieve the ID field from the dict.
	sID, ok, err := dictDefinition.Get(starlark.String("ID"))
	if err != nil || !ok {
		return fmt.Errorf("%w 'definition' getting ID field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}
	id, ok := sID.(starlark.String)
	if !ok {
		return fmt.Errorf("%w 'definition' ID field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}
	cd.ID = string(id)

	// Retrieve the Details field fro	m the dict.
	sName, ok, err := dictDefinition.Get(starlark.String("Details"))
	if err != nil || !ok {
		return fmt.Errorf("%w 'definition' getting name field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}
	name, ok := sName.(starlark.String)
	if !ok {
		return fmt.Errorf("%w 'definition' name field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}
	cd.Details = string(name)

	// Retrieve the BuildFunction field from the dict.
	sBuildFunction, ok, err := dictDefinition.Get(starlark.String("BuildFunction"))
	if err != nil || !ok {
		return fmt.Errorf("%w 'definition' getting build function field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}
	buildFunction, ok := sBuildFunction.(*starlark.Function)
	if !ok {
		return fmt.Errorf("%w 'definition' build function field must be a function %s", ErrCombinationDefinition, cd.StarScript)
	}
	cd.buildFunction = buildFunction

	return nil
}

func (cd *StarlarkDefinition) predeclared() starlark.StringDict {
	// uuidF() is a Go function called from Starlark.
	// It returns a new UUID.
	uuidF := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		uuidFunc := cd.getUuidFunc()
		id, err := uuidFunc()
		if err != nil {
			return nil, err

		}

		return starlark.String(id), nil
	}

	// now() is a Go function called from Starlark.
	// It returns the current time in RFC3339 format.
	now := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		now := cd.getNowFunc()()

		return starlark.String(time.Time.Format(now, time.RFC3339)), nil
	}

	// randomInt() is a Go function called from Starlark.
	// It returns multiple random values from an interval.
	randomInt := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var minVal uint
		var maxVal uint
		var nr int
		var allowDuplicates = false
		var sort = false

		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "min", &minVal, "max", &maxVal, "nr", &nr, "allow_duplicates?", &allowDuplicates, "sort?", &sort); err != nil {
			return nil, err
		}

		sliceResult := make([]uint, 0)

		for i := 0; i < nr; i++ {
			randUint, err := cd.getRandomIntFunc()(minVal, maxVal)
			if err != nil {
				return nil, err
			}

			if !allowDuplicates && slices.Contains(sliceResult, randUint) {
				i--
				continue
			}

			sliceResult = append(sliceResult, randUint)
		}

		if sort {
			slices.Sort(sliceResult)
		}

		result := starlark.NewList([]starlark.Value{})
		for _, randUint := range sliceResult {
			err := result.Append(starlark.MakeUint(randUint))
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	// renderTextTemplate() is a Go function called from Starlark.
	// It renders a text template with the given arguments.
	renderTextTemplate := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var tpl string
		var tplJsonArgs string
		var tplGoArgs any

		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "tpl", &tpl, "tplJsonArgs", &tplJsonArgs); err != nil {
			return nil, err
		}

		// Create a new template and parse the letter into it.
		t := template.Must(template.New("render_text_template").Parse(tpl))
		if err := json.Unmarshal([]byte(tplJsonArgs), &tplGoArgs); err != nil {
			return nil, fmt.Errorf("unmarshal slJson args: %w", err)
		}
		buff := &bytes.Buffer{}

		// Execute the template.
		err := t.Execute(buff, tplGoArgs)
		if err != nil {
			return nil, fmt.Errorf("execute template: %w", err)
		}

		return starlark.String(buff.String()), nil
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		// TODO: Move to module.
		"uuid":                 starlark.NewBuiltin("uuid", uuidF),
		"now":                  starlark.NewBuiltin("now", now),
		"random_int":           starlark.NewBuiltin("random_int", randomInt),
		"render_text_template": starlark.NewBuiltin("render_text_template", renderTextTemplate),
		"json":                 slJson.Module,
	}

	return predeclared
}

func (cd *StarlarkDefinition) getUuidFunc() func() (string, error) {
	if cd.uuidFunc != nil {
		return cd.uuidFunc
	}

	cd.uuidFunc = func() (string, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return id.String(), nil
	}

	return cd.uuidFunc
}

func (cd *StarlarkDefinition) getNowFunc() func() time.Time {
	if cd.nowFunc != nil {
		return cd.nowFunc
	}

	cd.nowFunc = func() time.Time {
		return time.Now()
	}

	return cd.nowFunc
}

func (cd *StarlarkDefinition) getRandomIntFunc() func(min uint, max uint) (uint, error) {
	if cd.randomUintFunc != nil {
		return cd.randomUintFunc
	}

	cd.randomUintFunc = random.NewCrypto().Uint

	return cd.randomUintFunc
}

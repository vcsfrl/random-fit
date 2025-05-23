package combination

import (
	"errors"
	"fmt"
	"github.com/vcsfrl/random-fit/internal/platform/starlark/random"
	"github.com/vcsfrl/random-fit/internal/platform/starlark/template"
	"github.com/vcsfrl/random-fit/internal/platform/starlark/uuid"
	slJson "go.starlark.net/lib/json"
	slTime "go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var ErrCombinationDefinition = errors.New("error combination definition")

type StarlarkDefinition struct {
	ID         string
	Details    string
	StarScript string

	buildFunction *starlark.Function
	thread        *starlark.Thread
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
	definition, hasDefinition := globals["definition"]
	if !hasDefinition {
		return fmt.Errorf("%w missing 'definition' Dict %s", ErrCombinationDefinition, cd.StarScript)
	}

	dictDefinition, hasDefinition := definition.(*starlark.Dict)
	if !hasDefinition {
		return fmt.Errorf("%w 'definition' must be a Dict %s", ErrCombinationDefinition, cd.StarScript)
	}

	// Retrieve the ID field from the dict.
	sID, hasDefinition, err := dictDefinition.Get(starlark.String("ID"))
	if err != nil || !hasDefinition {
		return fmt.Errorf("%w 'definition' getting ID field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}

	id, hasDefinition := sID.(starlark.String)
	if !hasDefinition {
		return fmt.Errorf("%w 'definition' ID field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}

	cd.ID = string(id)

	// Retrieve the Details field fro	m the dict.
	sName, hasDefinition, err := dictDefinition.Get(starlark.String("Details"))
	if err != nil || !hasDefinition {
		return fmt.Errorf("%w 'definition' getting name field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}

	name, hasDefinition := sName.(starlark.String)
	if !hasDefinition {
		return fmt.Errorf("%w 'definition' name field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}

	cd.Details = string(name)

	// Retrieve the BuildFunction field from the dict.
	sBuildFunction, hasDefinition, err := dictDefinition.Get(starlark.String("BuildFunction"))
	if err != nil || !hasDefinition {
		return fmt.Errorf("%w 'definition' getting build function field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}

	buildFunction, hasDefinition := sBuildFunction.(*starlark.Function)
	if !hasDefinition {
		return fmt.Errorf("%w 'definition' build function field must be a function %s", ErrCombinationDefinition, cd.StarScript)
	}

	cd.buildFunction = buildFunction

	return nil
}

func (cd *StarlarkDefinition) predeclared() starlark.StringDict {
	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"uuid":     uuid.Module,
		"template": template.Module,
		"json":     slJson.Module,
		"time":     slTime.Module,
		"random":   random.Module,
	}

	return predeclared
}

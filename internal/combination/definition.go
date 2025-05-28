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

type StarDefinition struct {
	ID         string
	Details    string
	StarScript string

	buildFunction *starlark.Function
	thread        *starlark.Thread

	UUIDModule     *uuid.UUID
	TemplateModule *template.Template
	RandomModule   *random.Random
}

func NewCombinationDefinition(script string) (*StarDefinition, error) {
	definition := &StarDefinition{
		StarScript: script,
	}

	err := definition.init()
	if err != nil {
		return nil, err
	}

	return definition, nil
}

func (cd *StarDefinition) CallScriptBuildFunction() (string, error) {
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

func (cd *StarDefinition) init() error {
	// The Thread defines the behavior of the built-in 'print' function.
	cd.thread = &starlark.Thread{
		Name: cd.Details,
		Print: func(_ *starlark.Thread, msg string) {
			fmt.Println(msg) //nolint:forbidigo
		},
	}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), cd.thread, cd.StarScript, nil, cd.predeclared())
	if err != nil {
		evalErr := &starlark.EvalError{}
		if errors.As(err, &evalErr) {
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

	if err := cd.initID(dictDefinition); err != nil {
		return err
	}

	if err := cd.initDetails(dictDefinition); err != nil {
		return err
	}

	if err := cd.initBuildFunction(dictDefinition); err != nil {
		return err
	}

	return nil
}

func (cd *StarDefinition) initBuildFunction(dictDefinition *starlark.Dict) error {
	// Retrieve the BuildFunction field from the dict.
	sBuildFunction, hasBuildFunction, err := dictDefinition.Get(starlark.String("BuildFunction"))
	if err != nil || !hasBuildFunction {
		return fmt.Errorf("%w 'definition' getting build function field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}

	buildFunction, hasBuildFunction := sBuildFunction.(*starlark.Function)
	if !hasBuildFunction {
		return fmt.Errorf(
			"%w 'definition' build function field must be a function %s",
			ErrCombinationDefinition,
			cd.StarScript,
		)
	}

	cd.buildFunction = buildFunction

	return nil
}

func (cd *StarDefinition) initDetails(dictDefinition *starlark.Dict) error {
	// Retrieve the Details field from the dict.
	sName, hasDetails, err := dictDefinition.Get(starlark.String("Details"))
	if err != nil || !hasDetails {
		return fmt.Errorf("%w 'definition' getting details field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}

	details, hasDetails := sName.(starlark.String)
	if !hasDetails {
		return fmt.Errorf("%w 'definition' details field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}

	cd.Details = string(details)

	return nil
}

func (cd *StarDefinition) initID(dictDefinition *starlark.Dict) error {
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

	return nil
}

func (cd *StarDefinition) predeclared() starlark.StringDict {
	cd.UUIDModule = uuid.New()
	cd.TemplateModule = template.New()
	cd.RandomModule = random.New()

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"uuid":     cd.UUIDModule.Module,
		"template": cd.TemplateModule.Module,
		"random":   cd.RandomModule.Module,
		"json":     slJson.Module,
		"time":     slTime.Module,
	}

	return predeclared
}

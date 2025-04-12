package combination

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/tmp/platform/random"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"time"
)

var ErrCombinationDefinition = fmt.Errorf("error combination definition")

type StarlarkDefinition struct {
	ID         string
	Name       string
	StarScript string
	GoTemplate string
	thread     *starlark.Thread

	buildFunction  *starlark.Function
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

func (cd *StarlarkDefinition) Generator() (func() (*Combination, error), error) {
	buildLambda, err := starlark.Call(cd.thread, cd.buildFunction, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination data: %w", ErrCombinationDefinition, err)
	}

	return func() (*Combination, error) {
		_, err := starlark.Call(cd.thread, buildLambda, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: error building combination data: %w", ErrCombinationDefinition, err)
		}

		uuid, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("%w: error building combination uuid: %w", ErrCombinationDefinition, err)
		}
		return &Combination{
			UUID:         uuid,
			DefinitionID: cd.ID,
			Name:         cd.Name,
			Template:     cd.GoTemplate,
			Data:         nil,
		}, nil
	}, nil
}

func (cd *StarlarkDefinition) init() error {
	// The Thread defines the behavior of the built-in 'print' function.
	cd.thread = &starlark.Thread{
		Name:  cd.Name,
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), cd.thread, cd.StarScript, nil, cd.predeclared())
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf("execution error: %w\n%s", evalErr, evalErr.Backtrace())
		}
		return fmt.Errorf("execution error: %w", err)
	}

	// Retrieve a module global.
	definition, ok := globals["definition"]
	if !ok {
		return fmt.Errorf("%w missing 'definition' dict %s", ErrCombinationDefinition, cd.StarScript)
	}

	dictDefinition, ok := definition.(*starlark.Dict)
	if !ok {
		return fmt.Errorf("%w 'definition' must be a Dict %s", ErrCombinationDefinition, cd.StarScript)
	}
	sName, ok, err := dictDefinition.Get(starlark.String("Name"))
	if err != nil || !ok {
		return fmt.Errorf("%w 'definition' getting name field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}
	name, ok := sName.(starlark.String)
	if !ok {
		return fmt.Errorf("%w 'definition' name field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}

	cd.Name = string(name)

	sBuildFunction, ok, err := dictDefinition.Get(starlark.String("BuildFunction"))
	if err != nil || !ok {
		return fmt.Errorf("%w 'definition' getting build function field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}
	buildFunction, ok := sBuildFunction.(*starlark.Function)
	if !ok {
		return fmt.Errorf("%w 'definition' build function field must be a function %s", ErrCombinationDefinition, cd.StarScript)
	}
	cd.buildFunction = buildFunction

	sGoTemplate, ok, err := dictDefinition.Get(starlark.String("GoTemplate"))
	if err != nil || !ok {
		return fmt.Errorf("%w 'definition' getting go template field %s: %w", ErrCombinationDefinition, cd.StarScript, err)
	}
	goTemplate, ok := sGoTemplate.(starlark.String)
	if !ok {
		return fmt.Errorf("%w 'definition' go template field must be a string %s", ErrCombinationDefinition, cd.StarScript)
	}
	cd.GoTemplate = goTemplate.String()

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
		var min uint
		var max uint
		var nr int
		var allowDuplicates = false

		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "min", &min, "max", &max, "nr", &nr, "allow_duplicates?", &allowDuplicates); err != nil {
			return nil, err
		}

		result := starlark.NewList([]starlark.Value{})
		uintFunc, err := cd.getRandomIntFunc()(min, max)
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

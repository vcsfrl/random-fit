package build

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/tmp/core"
	"github.com/vcsfrl/random-fit/internal/tmp/platform/random"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"time"
)

type StartCollectionBuilder struct {
	thread      *starlark.Thread
	builderFunc starlark.Value
	starFile    string

	uuidFunc       func() (string, error)
	nowFunc        func() time.Time
	randomUintFunc func(min uint, max uint) (uint, error)
}

func NewStartCollectionBuilder(starFile string) (*StartCollectionBuilder, error) {
	builder := &StartCollectionBuilder{starFile: starFile}
	err := builder.start()
	if err != nil {
		return nil, err
	}

	return builder, nil
}

func (s *StartCollectionBuilder) Build() (*core.Collection, error) {
	v, err := starlark.Call(s.thread, s.builderFunc, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("building collection from Starlark: %w", err)
	}

	collection := &core.Collection{}
	if err := json.Unmarshal([]byte(v.String()), collection); err != nil {
		return nil, fmt.Errorf("parsing collection from Starlark: %w", err)
	}

	return collection, nil
}

func (s *StartCollectionBuilder) start() error {

	// The Thread defines the behavior of the built-in 'print' function.
	s.thread = &starlark.Thread{
		Name:  "collection-builder",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), s.thread, s.starFile, nil, s.predeclared())
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf("execution error: %w\n%s", evalErr, evalErr.Backtrace())
		}
		return fmt.Errorf("execution error: %w", err)
	}

	// Retrieve a module global.
	buildCollection, ok := globals["build_collection"]
	if !ok {
		return fmt.Errorf("missing 'build_collection' function definition in %s", s.starFile)
	}

	s.builderFunc = buildCollection

	return nil
}

func (s *StartCollectionBuilder) predeclared() starlark.StringDict {
	// uuid() is a Go function called from Starlark.
	// It returns a new UUID.
	uuid := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		uuidFunc := s.getUuidFunc()
		id, err := uuidFunc()
		if err != nil {
			return nil, err

		}

		return starlark.String(id), nil
	}

	// now() is a Go function called from Starlark.
	// It returns the current time in RFC3339 format.
	now := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		now := s.getNowFunc()()

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

		for i := 0; i < nr; i++ {
			uintFunc, err := s.randomUintFunc(min, max)
			if err != nil {
				return nil, err
			}
			err = result.Append(starlark.MakeUint(uintFunc))
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"uuid":       starlark.NewBuiltin("uuid", uuid),
		"now":        starlark.NewBuiltin("now", now),
		"random_int": starlark.NewBuiltin("random_int", randomInt),
	}

	return predeclared
}

func (s *StartCollectionBuilder) getUuidFunc() func() (string, error) {
	if s.uuidFunc != nil {
		return s.uuidFunc
	}

	s.uuidFunc = func() (string, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return id.String(), nil
	}

	return s.uuidFunc
}

func (s *StartCollectionBuilder) getNowFunc() func() time.Time {
	if s.nowFunc != nil {
		return s.nowFunc
	}

	s.nowFunc = func() time.Time {
		return time.Now()
	}

	return s.nowFunc
}

func (s *StartCollectionBuilder) getRandomIntFunc() func(min uint, max uint) (uint, error) {
	if s.randomUintFunc != nil {
		return s.randomUintFunc
	}

	s.randomUintFunc = random.NewCrypto().Uint

	return s.randomUintFunc
}

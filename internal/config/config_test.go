package config

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/core"
	"go.starlark.net/starlark"
	"log"
	"testing"
	"time"
)

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

type ConfigSuite struct {
	suite.Suite

	testFolder string
}

func (suite *ConfigSuite) SetupTest() {
	suite.testFolder = "testdata/"
}

func (suite *ConfigSuite) TestFromScript() {

	builder := &StartCollectionBuilder{}
	builder.Start()
	collection := builder.Build()

	spew.Dump(collection)
}

type StartCollectionBuilder struct {
	thread      *starlark.Thread
	builderFunc starlark.Value
}

func (s *StartCollectionBuilder) Start() {
	// uuid() is a Go function called from Starlark.
	// It returns a new UUID version 7.
	uuid := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err

		}

		return starlark.String(id.String()), nil
	}

	// now() is a Go function called from Starlark.
	// It returns the current time in RFC3339 format.
	now := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		now := time.Now()

		return starlark.String(time.Time.Format(now, time.RFC3339)), nil
	}
	// The Thread defines the behavior of the built-in 'print' function.
	s.thread = &starlark.Thread{
		Name:  "collection-builder",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"uuid": starlark.NewBuiltin("uuid", uuid),
		"now":  starlark.NewBuiltin("now", now),
	}

	// Execute a program.
	globals, err := starlark.ExecFile(s.thread, "testdata/collection.star", nil, predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		log.Fatal(err)
	}

	// Retrieve a module global.
	buildCollection, ok := globals["build_collection"]
	if !ok {
		log.Fatal("build_collection not found")
	}

	s.builderFunc = buildCollection
}

func (s *StartCollectionBuilder) Build() core.Collection {
	v, err := starlark.Call(s.thread, s.builderFunc, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := core.Collection{}

	if err := json.Unmarshal([]byte(v.String()), &collection); err != nil {
		log.Fatal(err)
	}

	return collection
}

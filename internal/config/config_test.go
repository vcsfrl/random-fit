package config

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.starlark.net/starlark"
	"log"
	"strings"
	"testing"
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
	// repeat(str, n=1) is a Go function called from Starlark.
	// It behaves like the 'string * int' operation.
	repeat := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var s string
		var n int = 1
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "s", &s, "n?", &n); err != nil {
			return nil, err
		}
		return starlark.String(strings.Repeat(s, n)), nil
	}

	// The Thread defines the behavior of the built-in 'print' function.
	thread := &starlark.Thread{
		Name:  "example",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"greeting": starlark.String("hello"),
		"repeat":   starlark.NewBuiltin("repeat", repeat),
	}

	// Execute a program.
	globals, err := starlark.ExecFile(thread, "testdata/collection.star", nil, predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		log.Fatal(err)
	}

	// Retrieve a module global.
	buildCollection := globals["build_collection"]

	v, err := starlark.Call(thread, buildCollection, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Collection:", v)

}

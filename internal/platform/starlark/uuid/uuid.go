package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "uuid",
	Members: starlark.StringDict{
		"v7": starlark.NewBuiltin("v7", v7),
	},
}

var v7Func func() (string, error)

func v7(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	uuidFunc := getUuidFunc()
	id, err := uuidFunc()

	if err != nil {
		return nil, err

	}

	return starlark.String(id), nil
}

func getUuidFunc() func() (string, error) {
	if v7Func != nil {
		return v7Func
	}

	v7Func = func() (string, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return "", fmt.Errorf("error generating uuid: %v", err)
		}
		return id.String(), nil
	}

	return v7Func
}

func SetUuidFunc(f func() (string, error)) {
	v7Func = f
}

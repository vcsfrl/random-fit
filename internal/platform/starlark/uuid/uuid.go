package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type UUID struct {
	Module *starlarkstruct.Module
	v7Func func() (string, error)
}

func New() *UUID {
	uuidModule := &UUID{}
	uuidModule.init()

	return uuidModule
}

func (u *UUID) SetUUIDFunc(f func() (string, error)) {
	u.v7Func = f
}

func (u *UUID) init() {
	u.Module = &starlarkstruct.Module{
		Name: "uuid",
		Members: starlark.StringDict{
			"v7": starlark.NewBuiltin("v7", u.v7),
		},
	}
}

// v7 generates a UUIDv7 string. It is a wrapper around the uuid.NewV7 function.
//
//nolint:lll
func (u *UUID) v7(_ *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) { //nolint:ireturn
	uuidFunc := u.getUUIDFunc()
	uniqueID, err := uuidFunc()

	if err != nil {
		return nil, err
	}

	return starlark.String(uniqueID), nil
}

func (u *UUID) getUUIDFunc() func() (string, error) {
	if u.v7Func != nil {
		return u.v7Func
	}

	u.v7Func = func() (string, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return "", fmt.Errorf("error generating uuid: %w", err)
		}

		return id.String(), nil
	}

	return u.v7Func
}

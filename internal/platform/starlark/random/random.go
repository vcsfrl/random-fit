package random

import (
	"github.com/vcsfrl/random-fit/internal/platform/random"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"slices"
)

var Module = &starlarkstruct.Module{
	Name: "random",
	Members: starlark.StringDict{
		"uint": starlark.NewBuiltin("uint", getUint),
	},
}

var randomUintFunc func(min uint, max uint) (uint, error)

// getUint() is a Go function called from Starlark.
// It returns multiple random values from an interval of type uint.
func getUint(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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
		randUint, err := getUintFunc()(minVal, maxVal)
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

func getUintFunc() func(min uint, max uint) (uint, error) {
	if randomUintFunc != nil {
		return randomUintFunc
	}

	randomUintFunc = random.NewCrypto().Uint

	return randomUintFunc
}

func SetUintFunc(f func(min uint, max uint) (uint, error)) {
	randomUintFunc = f
}

package random

import (
	"fmt"
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

var randomUintFunc func(minValue uint, maxValue uint) (uint, error)

// getUint() is a Go function called from Starlark.
// It returns multiple random values from an interval of type uint.
func getUint(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var minVal, maxVal uint

	var number int

	var allowDuplicates, sort = false, false

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "min", &minVal, "max", &maxVal, "nr", &number, "allow_duplicates?", &allowDuplicates, "sort?", &sort); err != nil {
		return nil, fmt.Errorf("unpack args: %w", err)
	}

	sliceResult := make([]uint, 0)

	for index := 0; index < number; index++ {
		randUint, err := getUintFunc()(minVal, maxVal)
		if err != nil {
			return nil, err
		}

		if !allowDuplicates && slices.Contains(sliceResult, randUint) {
			index--

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
			return nil, fmt.Errorf("error generating random number: %v", err)
		}
	}

	return result, nil
}

func getUintFunc() func(minValue uint, maxValue uint) (uint, error) {
	if randomUintFunc != nil {
		return randomUintFunc
	}

	randomUintFunc = random.NewCrypto().Uint

	return randomUintFunc
}

func SetUintFunc(f func(minValue uint, maxValue uint) (uint, error)) {
	randomUintFunc = f
}

package random

import (
	"fmt"
	"github.com/vcsfrl/random-fit/internal/platform/random"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"slices"
)

type Random struct {
	Module         *starlarkstruct.Module
	randomUintFunc func(minValue uint, maxValue uint) (uint, error)
}

// New creates a new Starlark module for random number generation.
func New() *Random {
	randomModule := &Random{}
	randomModule.init()

	return randomModule
}

func (r *Random) SetUintFunc(f func(minValue uint, maxValue uint) (uint, error)) {
	r.randomUintFunc = f
}

func (r *Random) init() {
	r.Module = &starlarkstruct.Module{
		Name: "random",
		Members: starlark.StringDict{
			"uint": starlark.NewBuiltin("uint", r.getUint),
		},
	}
}

// getUint() is a Go function called from Starlark.
// It returns multiple random values from an interval of type uint.
//
//nolint:ireturn
func (r *Random) getUint(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var minVal, maxVal uint

	var number int

	var allowDuplicates, sort = false, false

	//nolint:lll
	if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "min", &minVal, "max", &maxVal, "nr", &number, "allow_duplicates?", &allowDuplicates, "sort?", &sort); err != nil {
		return nil, fmt.Errorf("unpack args: %w", err)
	}

	sliceResult := make([]uint, 0)

	for index := 0; index < number; index++ {
		randUint, err := r.getUintFunc()(minVal, maxVal)
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
			return nil, fmt.Errorf("error generating random number: %w", err)
		}
	}

	return result, nil
}

func (r *Random) getUintFunc() func(minValue uint, maxValue uint) (uint, error) {
	if r.randomUintFunc != nil {
		return r.randomUintFunc
	}

	r.randomUintFunc = random.NewCrypto().Uint

	return r.randomUintFunc
}

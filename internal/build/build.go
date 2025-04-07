package build

import (
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/model"
	"go.starlark.net/starlark"
)

type Builder struct {
}

func (b *Builder) Build(definition *model.Definition) *model.Combination {

	// Run the Starlark script from the definition to create a new combination.

	// Build the template from the definition.

	// Build the combination from the template and the data from the Starlark script.
	return &model.Combination{
		UUID:         uuid.New(),
		DefinitionId: definition.ID,
		Data:         starlark.NewDict(1),
		View:         "",
		ViewType:     "",
	}
}

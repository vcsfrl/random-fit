package combination

import (
	"github.com/google/uuid"
	"go.starlark.net/starlark"
)

type Combination struct {
	UUID         uuid.UUID
	DefinitionID string
	Name         string
	Data         map[string]any
	GoTemplate   string

	s starlark.Dict
}

package model

import (
	"github.com/google/uuid"
	"go.starlark.net/starlark"
)

type Combination struct {
	UUID         uuid.UUID
	DefinitionId string
	Data         starlark.Value
}

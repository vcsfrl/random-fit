package model

import (
	"bytes"
	"github.com/google/uuid"
	"go.starlark.net/starlark"
)

type Combination struct {
	UUID            uuid.UUID
	Definition      *CombinationDefinition
	Data            starlark.Value
	Output          *bytes.Buffer
	OutputExtension string
}

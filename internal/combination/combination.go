package combination

import (
	"github.com/google/uuid"
)

type Combination struct {
	UUID         uuid.UUID
	DefinitionID string
	Name         string
	Data         map[string]any
	GoTemplate   string
}

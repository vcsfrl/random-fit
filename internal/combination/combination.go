package combination

import (
	"github.com/google/uuid"
)

type Combination struct {
	UUID         uuid.UUID
	DefinitionID string
	Name         string
	GoTemplate   string
	Data         any
}

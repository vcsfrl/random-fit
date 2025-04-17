package combination

import (
	"github.com/google/uuid"
	"time"
)

type Combination struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Name         string
	GoTemplate   string
	JSONData     string
	GoData       any
}

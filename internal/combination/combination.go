package combination

import (
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/core"
	"time"
)

type Combination struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Name         string
	JSONData     string

	Template string
	Data     *core.Collection
}

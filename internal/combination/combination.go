package combination

import (
	"github.com/google/uuid"
)

type Combination struct {
	UUID       uuid.UUID
	Definition Definition
	Data       map[string]any
}

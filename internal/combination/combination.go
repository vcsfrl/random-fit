package combination

import (
	"bytes"
	"github.com/google/uuid"
	"time"
)

type Combination struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Name         string
	Template     string
	Data         *bytes.Buffer
	Output       []Output
}

type Output struct {
	Extension string
	MimeType  string
	Data      []byte
}

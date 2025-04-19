package combination

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Combination struct {
	UUID           uuid.UUID
	CreatedAt      time.Time
	DefinitionID   string
	DefinitionName string
	Data           map[string]*Data
}

type Data struct {
	Extension string
	MimeType  string
	Type      string
	Data      *bytes.Buffer
}

func (d *Data) UnmarshalJSON(data []byte) error {

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, value := range raw {
		switch key {
		case "Extension":
			d.Extension = string(value)
		case "MimeType":
			d.MimeType = string(value)
		case "Type":
			d.Type = string(value)
		case "Data":
			d.Data = bytes.NewBuffer(value)
		}
	}

	return nil

}

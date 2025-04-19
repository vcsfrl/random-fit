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
	Data           map[DataType]*Data
}

type DataType string

const DataTypeJson DataType = "json"
const DataTypeMd DataType = "markdown"
const DataTypeHtml DataType = "html"

var DataTypes = []DataType{
	DataTypeJson,
	DataTypeMd,
	DataTypeHtml,
}

type Data struct {
	Extension string
	MimeType  string
	Type      DataType
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
			dataType := DataType(value)
			// TODO: Check if dataType is valid
			//if !slices.Contains(DataTypes, dataType) {
			//	return fmt.Errorf("invalid data type: %s", dataType)
			//}
			d.Type = dataType
		case "Data":
			d.Data = bytes.NewBuffer(value)
		}
	}

	return nil

}

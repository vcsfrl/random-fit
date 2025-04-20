package combination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"slices"
	"strconv"
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
		unquoted, err := strconv.Unquote(string(value))
		if err != nil {
			return fmt.Errorf("error unquoting value: %w", err)
		}

		switch key {
		case "Extension":
			d.Extension = unquoted
		case "MimeType":
			d.MimeType = unquoted
		case "Type":
			dataType := DataType(unquoted)
			if !slices.Contains(DataTypes, dataType) {
				return fmt.Errorf("invalid data type: %s", dataType)
			}
			d.Type = dataType
		case "Data":
			d.Data = bytes.NewBuffer([]byte(unquoted))
		}
	}

	return nil

}

package combination

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"slices"
	"strconv"
	"time"
)

type Combination struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Details      string
	Data         map[DataType]*Data `validate:"combination_data_json"`
}

type DataType string

const DataTypeJson DataType = "json"
const DataTypeMd DataType = "markdown"

var DataTypes = []DataType{
	DataTypeJson,
	DataTypeMd,
}

type Data struct {
	Extension string
	MimeType  string
	Type      DataType
	Data      *bytes.Buffer
}

type gobData struct {
	Extension string
	MimeType  string
	Type      DataType
	Data      []byte
}

func (d *Data) GobEncode() ([]byte, error) {
	gData := &gobData{
		Extension: d.Extension,
		MimeType:  d.MimeType,
		Type:      d.Type,
		Data:      d.Data.Bytes(),
	}

	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(gData); err != nil {
		return nil, fmt.Errorf("error encoding data: %w", err)
	}
	return buffer.Bytes(), nil
}

func (d *Data) GobDecode(data []byte) error {
	gData := &gobData{}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(gData); err != nil {
		return fmt.Errorf("error decoding data: %w", err)
	}

	d.Extension = gData.Extension
	d.MimeType = gData.MimeType
	d.Type = gData.Type
	d.Data = bytes.NewBuffer(gData.Data)

	return nil
}

func (d *Data) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("error decoding data: %w", err)
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

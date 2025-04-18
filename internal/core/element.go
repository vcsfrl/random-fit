package core

import (
	"encoding/json"
	"fmt"
)

var ErrElement = fmt.Errorf("error element")

// Element is a generic type for an element.
type Element struct {
	Metadata Metadata
	Values   []Value
}

func (t *Element) UnmarshalJSON(data []byte) error {
	var element struct {
		Metadata Metadata
		Values   []any
	}

	if err := json.Unmarshal(data, &element); err != nil {
		return fmt.Errorf("%w: error unmarshalling: %w", ErrElement, err)
	}

	t.Metadata = element.Metadata

	for _, value := range element.Values {
		switch v := value.(type) {
		case string:
			t.Values = append(t.Values, StringValue{Value: v})
		case int:
			t.Values = append(t.Values, IntValue{Value: v})
		case float64:
			t.Values = append(t.Values, FloatValue{Value: v})
		case bool:
			t.Values = append(t.Values, BoolValue{Value: v})
		default:
			return fmt.Errorf("%w: unknown value type", ErrElement)
		}
	}

	return nil
}

// Value is an interface that represents a value.
type Value interface {
	String() string
	Type() string
}

// StringValue represents a string value.
type StringValue struct {
	Value string
}

func (v StringValue) String() string {
	return v.Value
}

func (v StringValue) Type() string {
	return "string"
}

// IntValue represents an integer value.
type IntValue struct {
	Value int
}

func (v IntValue) String() string {
	return fmt.Sprintf("%d", v.Value)
}

func (v IntValue) Type() string {
	return "int"
}

// FloatValue represents a float value.
type FloatValue struct {
	Value float64
}

func (v FloatValue) String() string {
	return fmt.Sprintf("%f", v.Value)
}

func (v FloatValue) Type() string {
	return "float"
}

// BoolValue represents a boolean value.
type BoolValue struct {
	Value bool
}

func (v BoolValue) String() string {
	return fmt.Sprintf("%t", v.Value)
}

func (v BoolValue) Type() string {
	return "bool"
}

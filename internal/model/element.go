package model

import (
	"fmt"
	"time"
)

// Element is a generic type for an element.
type Element struct {
	ID             string
	DefinitionId   string
	Name           string
	Values         []fmt.Stringer
	Date           time.Time
	ValueSeparator string
}

func (e Element) String() string {
	separator := " "
	if e.ValueSeparator != "" {
		separator = e.ValueSeparator
	}

	result := ""
	if len(e.Values) == 0 {
		return result
	}

	for i, v := range e.Values {
		if i > 0 {
			result += separator
		}
		result += v.String()
	}

	return result
}

// ElementValue is a generic type for the value of an element.
type ElementValue[T any] struct {
	Value T
}

func (e *ElementValue[T]) String() string {
	return fmt.Sprintf("%v", e.Value)
}

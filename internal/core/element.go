package core

import (
	"fmt"
)

// Element is a generic type for an element.
type Element struct {
	Metadata Metadata
	Values   []*ElementValue[any]
}

// ElementValue is a generic type for the value of an element.
type ElementValue[T any] struct {
	Value T
}

func (e *ElementValue[T]) String() string {
	return fmt.Sprintf("%v", e.Value)
}

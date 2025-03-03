package model

import (
	"fmt"
	"time"
)

// Element is a generic type for an element.
type Element struct {
	ID     string
	Name   string
	Values []fmt.Stringer
	Date   time.Time
}

func (e Element) String() string {
	return ""
}

// ElementValue is a generic type for the value of an element.
type ElementValue[T any] struct {
	Value T
}

func (e *ElementValue[T]) String() string {
	return fmt.Sprintf("%v", e.Value)
}

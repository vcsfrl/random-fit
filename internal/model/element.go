package model

import (
	"fmt"
	"time"
)

// Element is a generic type for an element.
type Element[T any] struct {
	ID    int
	Name  string
	Value ElementValue[T]
	Date  time.Time
}

func (e Element[T]) String() string {
	return e.Value.String()
}

// ElementValue is a generic type for the value of an element.
type ElementValue[T any] struct {
	Value T
}

func (e *ElementValue[T]) String() string {
	return fmt.Sprintf("%v", e.Value)
}

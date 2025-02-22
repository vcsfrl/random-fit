package model

import (
	"fmt"
	"time"
)

type Element struct {
	ID    int
	Name  string
	Value ElementValue[any]
	Date  time.Time
}

type ElementValue[T comparable] struct {
	Value T
}

func (e *ElementValue[T]) String() string {
	return fmt.Sprintf("%v", e.Value)
}

func (e Element) String() string {
	return e.Value.String()
}

type ElementDefinition[T any] struct {
	ID          int
	Name        string
	Description string
	Options     T
	NrOfPicks   int
	UniquePicks bool
	GlobalPicks bool
}

type ElementOptionsInterval[T comparable] struct {
	Min T
	Max T
}

type ElementOptionsValues struct {
	Values []ElementValue[any]
}

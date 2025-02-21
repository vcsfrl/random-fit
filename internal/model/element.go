package model

import (
	"fmt"
	"time"
)

type Element struct {
	ID    int
	Name  string
	Value fmt.Stringer
	Date  time.Time
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

type ElementOptionsInterval struct {
	Min int
	Max int
}

type ElementOptionsValues struct {
	Values []fmt.Stringer
}

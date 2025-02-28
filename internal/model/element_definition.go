package model

// ElementDefinition is a generic type for the definition of an element.
type ElementDefinition[T any] struct {
	ID          int
	Name        string
	Description string
	NrOfPicks   int
	UniquePicks bool
	GlobalPicks bool

	Options ElementDefinitionOptions[T]
}

type ElementDefinitionOptions[T any] struct {
	Type     string
	Interval ElementDefinitionOptionsInterval[T]
	Values   []T
}

type ElementDefinitionOptionsInterval[T any] struct {
	MinValue T
	MaxValue T
}

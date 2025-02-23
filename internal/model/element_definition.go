package model

// ElementDefinition is a generic type for the definition of an element.
type ElementDefinition[T any] struct {
	ID          int
	Name        string
	Description string
	Options     T
	NrOfPicks   int
	UniquePicks bool
	GlobalPicks bool
}

// ElementOptionsInterval is a generic type for the options of an element with an interval.
// Values of an element will be generated within the interval.
type ElementOptionsInterval[T comparable] struct {
	Min ElementValue[T]
	Max ElementValue[T]
}

// ElementOptionsValues is a generic type for the options of an element with values.
// Values of an element will be picked from the given values.
type ElementOptionsValues struct {
	Values []ElementValue[any]
}

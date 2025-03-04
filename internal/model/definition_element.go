package model

// ElementDefinition is a generic type for the definition of an element.
type ElementDefinition struct {
	Identity     DefinitionIdentity
	UniquePicks  bool
	NrOfPicks    int
	PickStrategy PickStrategy
	Options      ElementDefinitionOptions
}

// ElementDefinitionOptions is a generic type for the options of an element definition.
type ElementDefinitionOptions struct {
	Interval ElementDefinitionOptionInterval[any]
	Values   []any
}

// ElementDefinitionOptionInterval is a generic type for the interval option of an element definition.
type ElementDefinitionOptionInterval[T any] struct {
	MinValue T
	MaxValue T
}

package model

type PickStrategy string

const PickStrategyRandom PickStrategy = "random"
const PickStrategyDice PickStrategy = "dice"

// ElementDefinition is a generic type for the definition of an element.
type ElementDefinition struct {
	ID           string
	Name         string
	Description  string
	UniquePicks  bool
	NrOfPicks    int
	PickStrategy PickStrategy
	Options      *ElementDefinitionOptions
}

type ElementDefinitionOptions struct {
	Interval *ElementDefinitionOptionInterval[any]
	Values   []any
}
type ElementDefinitionOptionInterval[T any] struct {
	MinValue T
	MaxValue T
}

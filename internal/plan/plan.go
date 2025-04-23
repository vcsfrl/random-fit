package plan

import (
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/combination"
	"time"
)

// Definition is used to define how to generate Combinations.
// It is used to generate combinations and  categorise them.
type Definition struct {
	ID      string
	Details string
	Users   []string
	GroupDefinition
}

// GroupDefinition is used to define a group of combinations.
type GroupDefinition struct {
	NamePrefix       string
	NumberOfGroups   int
	NrOfCombinations int
}

type Group struct {
	Details      string
	Combinations []*combination.Combination
}

type Plan struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Details      string
	Users        []string
	Groups       []*Group
}

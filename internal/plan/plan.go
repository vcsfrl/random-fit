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
	UserData
}

// UserData is used to define a group of combinations.
type UserData struct {
	ContainerName            string
	RecurrentGroupNamePrefix string
	RecurrentGroups          int
	NrOfGroupCombinations    int
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
	UserGroups   map[string][]*Group
}

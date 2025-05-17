package plan

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/combination"
	"os"
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
	ContainerName            []string
	RecurrentGroupNamePrefix string
	RecurrentGroups          int
	NrOfGroupCombinations    int
}

type GroupDetails struct {
	ContainerName []string
	Details       string
	User          string
}
type Group struct {
	GroupDetails
	Combinations []*combination.Combination
}

type PlanDetails struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Details      string
}

type Plan struct {
	PlanDetails
	UserGroups map[string][]*Group
}

type PlanCombination struct {
	PlanDetails
	GroupDetails
	Combination *combination.Combination
	Err         error
}

func NewJsonDefinition(fileName string) (*Definition, error) {
	result := &Definition{}

	// Read the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read json definition: %w", err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json definition: %w", err)
	}

	return result, nil
}

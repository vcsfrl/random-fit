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
	ID      string   `json:"id"`
	Details string   `json:"details"`
	Users   []string `json:"users"`
	UserData
}

// UserData is used to define a group of combinations.
type UserData struct {
	ContainerName            []string `json:"containerName"`
	RecurrentGroupNamePrefix string   `json:"recurrentGroupNamePrefix"`
	RecurrentGroups          int      `json:"recurrentGroups"`
	NrOfGroupCombinations    int      `json:"nrOfGroupCombinations"`
}

type Group struct {
	ContainerName []string
	Details       string
	User          string
}
type GroupCombination struct {
	Group
	Combinations []*combination.Combination
}

type Plan struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Details      string
}

type UserPlan struct {
	Plan
	UserGroups map[string][]*GroupCombination
}

type PlannedCombination struct {
	Plan
	Group
	Combination   *combination.Combination
	GroupSerialID int
	Err           error
}

func NewJSONDefinition(fileName string) (*Definition, error) {
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

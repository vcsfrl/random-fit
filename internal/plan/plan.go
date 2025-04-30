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
	ContainerName            string
	RecurrentGroupNamePrefix string
	RecurrentGroups          int
	NrOfGroupCombinations    int
}

type Group struct {
	ContainerName string
	Details       string
	Combinations  []*combination.Combination
}

type Plan struct {
	UUID         uuid.UUID
	CreatedAt    time.Time
	DefinitionID string
	Details      string
	UserGroups   map[string][]*Group
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

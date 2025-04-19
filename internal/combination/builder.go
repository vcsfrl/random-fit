package combination

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Builder interface {
	Build() (*Combination, error)
}

type StarlarkBuilder struct {
	definition *StarlarkDefinition
}

func NewStarlarkBuilder(definition *StarlarkDefinition) *StarlarkBuilder {
	return &StarlarkBuilder{definition: definition}
}

func (s *StarlarkBuilder) Build() (*Combination, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination uuid: %w", ErrCombinationDefinition, err)
	}

	combinationData, err := s.definition.CallScriptBuildFunction()
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination data: %w", ErrCombinationDefinition, err)
	}

	result := &Combination{
		UUID:           uuidV7,
		CreatedAt:      time.Now(),
		DefinitionID:   s.definition.ID,
		DefinitionName: s.definition.Name,
		Data:           make(map[DataType]*Data),
	}

	err = json.Unmarshal([]byte(combinationData), &result.Data)
	if err != nil {
		return nil, fmt.Errorf("%w: error unmarshalling combination data: %w", ErrCombinationDefinition, err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("%w: combination data is nil", ErrCombinationDefinition)
	}

	// Check if the Data map has a json key

	if _, ok := result.Data[DataTypeJson]; !ok {
		return nil, fmt.Errorf("%w: combination data does not contain json key", ErrCombinationDefinition)
	}

	return result, nil
}

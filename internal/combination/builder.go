package combination

import (
	"bytes"
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

	return &Combination{
		UUID:           uuidV7,
		CreatedAt:      time.Now(),
		DefinitionID:   s.definition.ID,
		DefinitionName: s.definition.Name,
		Data:           bytes.NewBuffer([]byte(combinationData)),
	}, nil
}

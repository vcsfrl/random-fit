package combination

import (
	"fmt"
	"github.com/google/uuid"
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
		UUID:         uuidV7,
		DefinitionID: s.definition.ID,
		Name:         s.definition.Name,
		GoTemplate:   s.definition.GoTemplate,
		Data:         combinationData,
	}, nil
}

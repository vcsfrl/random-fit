package combination

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Builder interface {
	Build() (*Combination, error)
}

type StarlarkBuilder struct {
	definition *StarlarkDefinition
	now        func() time.Time
	uuidV7     func() (uuid.UUID, error)
}

func NewStarlarkBuilder(definition *StarlarkDefinition) *StarlarkBuilder {
	return &StarlarkBuilder{
		definition: definition,
		now:        time.Now,
		uuidV7: func() (uuid.UUID, error) {
			return uuid.NewV7()
		},
	}
}

func (s *StarlarkBuilder) Build() (*Combination, error) {
	uuidV7, err := s.uuidV7()
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination uuid: %w", ErrCombinationDefinition, err)
	}

	combinationData, err := s.definition.CallScriptBuildFunction()
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination data: %w", ErrCombinationDefinition, err)
	}

	result := &Combination{
		UUID:         uuidV7,
		CreatedAt:    s.now(),
		DefinitionID: s.definition.ID,
		Details:      s.definition.Details,
		Data:         make(map[DataType]*Data),
	}

	err = json.Unmarshal([]byte(combinationData), &result.Data)
	if err != nil {
		return nil, fmt.Errorf("%w: error unmarshalling combination data: %w", ErrCombinationDefinition, err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("%w: combination data is nil", ErrCombinationDefinition)
	}

	validate, err := Validator()
	if err != nil {
		return nil, fmt.Errorf("%w: error creating validator: %w", ErrCombinationDefinition, err)
	}

	err = validate.Struct(result)
	if err != nil {
		return nil, fmt.Errorf("%w: error validating combination data: %w", ErrCombinationDefinition, err.(validator.ValidationErrors))
	}

	return result, nil
}

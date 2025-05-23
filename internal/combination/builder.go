package combination

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

var ErrStarBuilder = errors.New("error combination definition")

type Builder interface {
	Build() (*Combination, error)
}

type StarBuilder struct {
	definition *StarlarkDefinition
	now        func() time.Time
	uuidV7     func() (uuid.UUID, error)
	validate   *validator.Validate
}

func NewStarBuilder(definition *StarlarkDefinition) (*StarBuilder, error) {
	validate, err := Validator()
	if err != nil {
		return nil, fmt.Errorf("%w: error creating validator: %w", ErrStarBuilder, err)
	}

	return &StarBuilder{
		definition: definition,
		now:        time.Now,
		uuidV7: func() (uuid.UUID, error) {
			return uuid.NewV7()
		},
		validate: validate,
	}, nil
}

func (s *StarBuilder) Build() (*Combination, error) {
	uuidV7, err := s.uuidV7()
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination uuid: %w", ErrStarBuilder, err)
	}

	combinationData, err := s.definition.CallScriptBuildFunction()
	if err != nil {
		return nil, fmt.Errorf("%w: error building combination data: %w", ErrStarBuilder, err)
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
		return nil, fmt.Errorf("%w: error unmarshalling combination data: %w", ErrStarBuilder, err)
	}

	if result.Data == nil {
		return nil, fmt.Errorf("%w: combination data is nil", ErrStarBuilder)
	}

	err = s.validate.Struct(result)
	if err != nil {
		return nil, fmt.Errorf("%w: error validating combination data: %w", ErrStarBuilder, func() validator.ValidationErrors {
			var target validator.ValidationErrors
			_ = errors.As(err, &target)
			return target
		}())
	}

	return result, nil
}

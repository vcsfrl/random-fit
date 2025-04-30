package plan

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/combination"
	"time"
)

var ErrPlanBuild = fmt.Errorf("error building plan")

type Builder struct {
	Definition         *Definition
	Now                func() time.Time
	UuidV7             func() (uuid.UUID, error)
	CombinationBuilder combination.Builder
}

func NewBuilderFromStarConfig(combinationFile string, planFile string) *Builder {
	combinationDefinition, err := combination.NewCombinationDefinition(combinationFile)
	if err != nil {
		panic(fmt.Errorf("%w: error creating combination definition: %w", ErrPlanBuild, err))
	}

	planDefinition, err := NewJsonDefinition(planFile)
	if err != nil {
		panic(fmt.Errorf("%w: error creating plan definition: %w", ErrPlanBuild, err))
	}

	builder := combination.NewStarlarkBuilder(combinationDefinition)

	return &Builder{
		Definition:         planDefinition,
		Now:                time.Now,
		UuidV7:             uuid.NewV7,
		CombinationBuilder: builder,
	}

}

func NewBuilder(definition *Definition, builder combination.Builder) *Builder {
	return &Builder{
		Definition:         definition,
		Now:                time.Now,
		UuidV7:             uuid.NewV7,
		CombinationBuilder: builder,
	}
}

func (b *Builder) Build() (*Plan, error) {
	uuidV7, err := b.UuidV7()
	if err != nil {
		return nil, err
	}

	plan := &Plan{
		UUID:         uuidV7,
		CreatedAt:    b.Now(),
		DefinitionID: b.Definition.ID,
		Details:      b.Definition.Details,
		UserGroups:   make(map[string][]*Group),
	}

	for _, user := range b.Definition.Users {
		userGroups := make([]*Group, 0)

		// Create groups
		for i := 0; i < b.Definition.RecurrentGroups; i++ {
			group := &Group{
				Details:       fmt.Sprintf("%s-%d", b.Definition.RecurrentGroupNamePrefix, i+1),
				Combinations:  make([]*combination.Combination, 0),
				ContainerName: b.Definition.ContainerName,
			}
			for j := 0; j < b.Definition.NrOfGroupCombinations; j++ {
				newCombination, err := b.CombinationBuilder.Build()
				if err != nil {
					return nil, fmt.Errorf("%w: error building combination: %w", ErrPlanBuild, err)
				}

				group.Combinations = append(group.Combinations, newCombination)
			}

			userGroups = append(userGroups, group)
		}

		plan.UserGroups[user] = userGroups
	}

	return plan, nil
}

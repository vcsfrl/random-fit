package plan

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vcsfrl/random-fit/internal/combination"
	"time"
)

var ErrPlanBuild = errors.New("error building plan")
var ErrPlanBuildTerminated = fmt.Errorf("%w: build terminated", ErrPlanBuild)

const GeneratorBuffer = 1000

type Builder struct {
	Definition         *Definition
	Now                func() time.Time
	UUIDV7             func() (uuid.UUID, error)
	CombinationBuilder combination.Builder
}

func NewBuilderFromStarConfig(combinationFile string, planFile string) *Builder {
	combinationDefinition, err := combination.NewCombinationDefinition(combinationFile)
	if err != nil {
		panic(fmt.Errorf("%w: error creating combination definition: %w", ErrPlanBuild, err))
	}

	planDefinition, err := NewJSONDefinition(planFile)
	if err != nil {
		panic(fmt.Errorf("%w: error creating plan definition: %w", ErrPlanBuild, err))
	}

	builder, err := combination.NewStarBuilder(combinationDefinition)
	if err != nil {
		panic(fmt.Errorf("%w: error creating combination builder: %w", ErrPlanBuild, err))
	}

	return &Builder{
		Definition:         planDefinition,
		Now:                time.Now,
		UUIDV7:             uuid.NewV7,
		CombinationBuilder: builder,
	}
}

func NewBuilder(definition *Definition, builder combination.Builder) *Builder {
	return &Builder{
		Definition:         definition,
		Now:                time.Now,
		UUIDV7:             uuid.NewV7,
		CombinationBuilder: builder,
	}
}

func (b *Builder) Build() (*UserPlan, error) {
	uuidV7, err := b.UUIDV7()
	if err != nil {
		return nil, err
	}

	plan := &UserPlan{
		Plan: Plan{
			UUID:         uuidV7,
			CreatedAt:    b.Now(),
			DefinitionID: b.Definition.ID,
			Details:      b.Definition.Details,
		},
		UserGroups: make(map[string][]*GroupCombination),
	}

	for _, user := range b.Definition.Users {
		userGroups := make([]*GroupCombination, 0)

		// Create groups
		for i := range b.Definition.RecurrentGroups {
			group := &GroupCombination{
				Group: Group{
					Details:       fmt.Sprintf("%s-%d", b.Definition.RecurrentGroupNamePrefix, i+1),
					ContainerName: b.Definition.ContainerName,
					User:          user,
				},
				Combinations: make([]*combination.Combination, 0),
			}

			for range b.Definition.NrOfGroupCombinations {
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

func (b *Builder) Generate(ctx context.Context) chan *PlannedCombination {
	generator := make(chan *PlannedCombination, GeneratorBuffer)

	uuidV7, err := b.UUIDV7()
	if err != nil {
		generator <- &PlannedCombination{Err: fmt.Errorf("%w: error creating uuid v7: %w", ErrPlanBuild, err)}
		close(generator)

		return generator
	}

	createdAt := b.Now()

	go func() {
		defer close(generator)

		for _, user := range b.Definition.Users {
			// Create groups
			for i := range b.Definition.RecurrentGroups {
				for j := range b.Definition.NrOfGroupCombinations {
					select {
					case <-ctx.Done():
						generator <- &PlannedCombination{Err: ErrPlanBuildTerminated}

						return
					default: // continue
					}

					newCombination, err := b.CombinationBuilder.Build()
					plannedCombination := &PlannedCombination{
						Plan: Plan{
							UUID:         uuidV7,
							CreatedAt:    createdAt,
							DefinitionID: b.Definition.ID,
							Details:      b.Definition.Details,
						},
						Group: Group{
							Details:       fmt.Sprintf("%s-%d", b.Definition.RecurrentGroupNamePrefix, i+1),
							ContainerName: b.Definition.ContainerName,
							User:          user,
						},
						Combination:   newCombination,
						GroupSerialID: j + 1,
						Err:           err,
					}

					generator <- plannedCombination
				}
			}
		}
	}()

	return generator
}

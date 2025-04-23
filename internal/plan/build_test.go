package plan

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/combination"
	"testing"
	"time"
)

func TestBuildSuite(t *testing.T) {
	suite.Run(t, new(BuildSuite))
}

type BuildSuite struct {
	suite.Suite
}

func (suite *BuildSuite) TestBuild() {
	definition := &Definition{
		ID:      "test",
		Details: "Test",
		GroupDefinition: GroupDefinition{
			NamePrefix:       "Test",
			NumberOfGroups:   4,
			NrOfCombinations: 3,
		},
	}

	// Mock the combination builder
	mockBuilder := &MockCombinationBuilder{}
	plan, err := NewBuilder(definition, mockBuilder).Build()
	suite.NoError(err)
	suite.NotNil(plan)

	suite.Equal(definition.ID, plan.DefinitionID)
	suite.Equal(definition.Details, plan.Details)
	suite.Equal(definition.GroupDefinition.NumberOfGroups, len(plan.Groups))
	suite.Equal(definition.GroupDefinition.NrOfCombinations, len(plan.Groups[0].Combinations))
	suite.Equal(definition.GroupDefinition.NamePrefix+"-1", plan.Groups[0].Details)
	suite.Equal(definition.GroupDefinition.NamePrefix+"-2", plan.Groups[1].Details)
	suite.Equal(definition.GroupDefinition.NamePrefix+"-3", plan.Groups[2].Details)
	suite.Equal(definition.GroupDefinition.NamePrefix+"-4", plan.Groups[3].Details)
	suite.Equal(mockBuilder.Calls, definition.GroupDefinition.NumberOfGroups*definition.GroupDefinition.NrOfCombinations)
	suite.Equal(mockBuilder.Calls, len(plan.Groups[0].Combinations)+len(plan.Groups[1].Combinations)+len(plan.Groups[2].Combinations)+len(plan.Groups[3].Combinations))

	suite.Equal("test-1", plan.Groups[0].Combinations[0].Details)
	suite.Equal("test-12", plan.Groups[3].Combinations[2].Details)
}

type MockCombinationBuilder struct {
	Calls int
}

func (m *MockCombinationBuilder) Build() (*combination.Combination, error) {
	m.Calls++
	return &combination.Combination{
		UUID:         uuid.New(),
		CreatedAt:    time.Now(),
		DefinitionID: "test",
		Details:      fmt.Sprintf("test-%d", m.Calls),
		Data:         make(map[combination.DataType]*combination.Data),
	}, nil
}

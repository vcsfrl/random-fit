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
		Users:   []string{"user-1", "user-2"},
		UserData: UserData{
			RecurrentGroupNamePrefix: "Test",
			RecurrentGroups:          4,
			NrOfGroupCombinations:    3,
		},
	}

	// Mock the combination builder
	mockBuilder := &MockCombinationBuilder{}
	plan, err := NewBuilder(definition, mockBuilder).Build()
	suite.NoError(err)
	suite.NotNil(plan)

	suite.Equal(definition.ID, plan.DefinitionID)
	suite.Equal(definition.Details, plan.Details)
	suite.Equal(definition.RecurrentGroups, len(plan.UserGroups["user-1"]))
	suite.Equal(definition.NrOfGroupCombinations, len(plan.UserGroups["user-1"][0].Combinations))
	suite.Equal(definition.RecurrentGroupNamePrefix+"-1", plan.UserGroups["user-1"][0].Details)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-2", plan.UserGroups["user-1"][1].Details)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-3", plan.UserGroups["user-1"][2].Details)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-4", plan.UserGroups["user-1"][3].Details)
	suite.Equal("test-1", plan.UserGroups["user-1"][0].Combinations[0].Details)
	suite.Equal("test-12", plan.UserGroups["user-1"][3].Combinations[2].Details)

	suite.Equal(definition.RecurrentGroups, len(plan.UserGroups["user-2"]))
	suite.Equal(definition.NrOfGroupCombinations, len(plan.UserGroups["user-2"][0].Combinations))
	suite.Equal(definition.RecurrentGroupNamePrefix+"-1", plan.UserGroups["user-2"][0].Details)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-2", plan.UserGroups["user-2"][1].Details)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-3", plan.UserGroups["user-2"][2].Details)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-4", plan.UserGroups["user-2"][3].Details)
	suite.Equal("test-13", plan.UserGroups["user-2"][0].Combinations[0].Details)
	suite.Equal("test-24", plan.UserGroups["user-2"][3].Combinations[2].Details)

	suite.Equal(mockBuilder.Calls, definition.RecurrentGroups*definition.NrOfGroupCombinations*len(definition.Users))
}

func (suite *BuildSuite) TestGenerate() {
	definition := &Definition{
		ID:      "test",
		Details: "Test",
		Users:   []string{"user-1", "user-2"},
		UserData: UserData{
			ContainerName:            []string{"test1"},
			RecurrentGroupNamePrefix: "Test",
			RecurrentGroups:          4,
			NrOfGroupCombinations:    3,
		},
	}

	// Mock the combination builder
	mockBuilder := &MockCombinationBuilder{}
	generator := NewBuilder(definition, mockBuilder).Generate()
	suite.NotNil(generator)

	data := []PlannedCombination{}
	for genCombination := range generator {
		data = append(data, *genCombination)
	}

	suite.Equal(definition.RecurrentGroups*definition.NrOfGroupCombinations*len(definition.Users), len(data))
	suite.Equal(definition.RecurrentGroups*definition.NrOfGroupCombinations*len(definition.Users), mockBuilder.Calls)

	suite.NoError(data[0].Err)
	suite.Equal(definition.ID, data[0].Plan.DefinitionID)
	suite.Equal(definition.Details, data[0].Plan.Details)
	suite.NotNil(data[0].Plan.UUID)
	suite.NotNil(data[0].Plan.CreatedAt)
	suite.Equal(definition.RecurrentGroupNamePrefix+"-1", data[0].Group.Details)
	suite.Equal(definition.UserData.ContainerName, data[0].Group.ContainerName)
	suite.Equal(definition.Users[0], data[0].Group.User)
	suite.Equal(1, data[0].GroupSerialId)
	suite.NotNil(data[0].Combination)
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

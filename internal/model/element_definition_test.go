package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementDefinition(t *testing.T) {
	definitionIntervalInt := ElementDefinition{
		ID:          "test-definition-1",
		Name:        "Test",
		Description: "Test description",
		NrOfPicks:   1,
		UniquePicks: true,
		Options: &ElementDefinitionOptions{
			Interval: &ElementDefinitionOptionInterval[any]{
				MinValue: 0,
				MaxValue: 10,
			},
			Values: []any{"Test 1", "Test 2"},
		},
	}

	// Dummy test to avoid compilation error
	assert.NotNil(t, definitionIntervalInt)
	//assert.NotEqual(t, definitionIntervalInt, definitionIntervalString)
	//assert.NotEqual(t, definitionIntervalInt, definitionValues)
}

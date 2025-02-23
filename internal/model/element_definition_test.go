package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementDefinition(t *testing.T) {
	definitionIntervalInt := ElementDefinition[ElementOptionsInterval[int]]{
		ID:          1,
		Name:        "Test",
		Description: "Test description",
		Options:     ElementOptionsInterval[int]{Min: ElementValue[int]{Value: 1}, Max: ElementValue[int]{Value: 100}},
		NrOfPicks:   1,
		UniquePicks: true,
		GlobalPicks: false,
	}

	definitionIntervalString := ElementDefinition[ElementOptionsInterval[string]]{
		ID:          2,
		Name:        "Test",
		Description: "Test description",
		Options:     ElementOptionsInterval[string]{Min: ElementValue[string]{Value: "A"}, Max: ElementValue[string]{Value: "Z"}},
		NrOfPicks:   1,
		UniquePicks: true,
		GlobalPicks: false,
	}

	definitionValues := ElementDefinition[ElementOptionsValues]{
		ID:          3,
		Name:        "Test",
		Description: "Test description",
		Options:     ElementOptionsValues{Values: []ElementValue[any]{{Value: 1}, {Value: "Test"}}},
		NrOfPicks:   1,
		UniquePicks: true,
		GlobalPicks: false,
	}

	// Dummy test to avoid compilation error
	assert.NotEqual(t, definitionIntervalInt, definitionIntervalString)
	assert.NotEqual(t, definitionIntervalInt, definitionValues)
}

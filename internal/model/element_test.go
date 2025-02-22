package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementDefinition(t *testing.T) {
	definitionInterval := ElementDefinition[ElementOptionsInterval]{
		ID:          1,
		Name:        "Test",
		Description: "Test description",
		Options:     ElementOptionsInterval{Min: 1, Max: 10},
		NrOfPicks:   1,
		UniquePicks: true,
		GlobalPicks: false,
	}

	definitionValues := ElementDefinition[ElementOptionsValues]{
		ID:          2,
		Name:        "Test",
		Description: "Test description",
		Options:     ElementOptionsValues{Values: []ElementValue[any]{{Value: 1}, {Value: "Test"}}},
		NrOfPicks:   1,
		UniquePicks: true,
		GlobalPicks: false,
	}

	assert.NotEqual(t, definitionInterval, definitionValues)
}

func TestElementValue_String(t *testing.T) {
	intValue := ElementValue[int]{Value: 1}
	assert.Equal(t, "1", intValue.String())

	stringValue := ElementValue[string]{Value: "Test"}
	assert.Equal(t, "Test", stringValue.String())
}

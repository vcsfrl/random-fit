package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementDefinition(t *testing.T) {
	definitionIntervalInt := ElementDefinition[int]{
		ID:          1,
		Name:        "Test",
		Description: "Test description",
		NrOfPicks:   1,
		UniquePicks: true,
		GlobalPicks: false,
		Options: ElementDefinitionOptions[int]{
			Type: "interval",
			Interval: ElementDefinitionOptionsInterval[int]{
				MinValue: 0,
				MaxValue: 10,
			},
			Values: nil,
		},
	}

	// Dummy test to avoid compilation error
	assert.NotNil(t, definitionIntervalInt)
	//assert.NotEqual(t, definitionIntervalInt, definitionIntervalString)
	//assert.NotEqual(t, definitionIntervalInt, definitionValues)
}

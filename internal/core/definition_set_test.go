package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetDefinition(t *testing.T) {
	setDefinition := SetDefinition{
		Metadata: DefinitionMetadata{
			ID:          "test-definition-1",
			Name:        "Test",
			Description: "Test description",
		},
		Elements: []ElementDefinition{
			{
				Metadata: DefinitionMetadata{
					ID:          "element-1",
					Name:        "Element 1",
					Description: "Description of the element",
				},
				UniquePicks:  false,
				NrOfPicks:    2,
				PickStrategy: PickStrategyRandom,
				Options: ElementDefinitionOptions{
					Interval: ElementDefinitionOptionInterval[any]{
						MinValue: 0,
						MaxValue: 10,
					},
					Values: []any{"Test 1", "Test 2"},
				},
			},
		},
	}

	assert.NotNil(t, setDefinition)
}

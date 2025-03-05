package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectionDefinition(t *testing.T) {
	collectionDefinition := CollectionDefinition{
		Identity: DefinitionIdentity{
			ID:          "collection-definition-1",
			Name:        "Test",
			Description: "Test description",
		},
		Sets: []SetDefinition{
			{
				Identity: DefinitionIdentity{
					ID:          "set-definition-1",
					Name:        "Test",
					Description: "Test description",
				},
				Elements: []ElementDefinition{
					{
						Identity: DefinitionIdentity{
							ID:          "element-definition-1",
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
			},
		},
	}

	assert.NotNil(t, collectionDefinition)

}

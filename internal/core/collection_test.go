package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	collection := Collection{
		Identity: Identity{
			ID:           "collection-1",
			DefinitionID: "definition-1",
			Name:         "Test",
			Description:  "Description of the collection",
			Date:         time.Time{},
		},
		Sets: []Set{
			{
				Identity: Identity{
					ID:          "Set 1",
					Name:        "Test",
					Description: "Description of the set",
					Date:        time.Time{},
				},
				Elements: []Element{
					{
						Identity: Identity{
							ID:           "element-1",
							DefinitionID: "definition-1",
							Name:         "Element 1",
							Description:  "Description of the element",
							Date:         time.Time{},
						},
					},
				},
			},
		},
	}

	assert.NotNil(t, collection)

}

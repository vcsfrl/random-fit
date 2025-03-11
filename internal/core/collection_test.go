package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	collection := Collection{
		Metadata: Metadata{
			ID:          "collection-1",
			Name:        "Test Collection 1",
			Description: "Description of the collection",
			Date:        time.Time{},
		},
		Sets: []*Set{
			{
				Metadata: Metadata{
					ID:          "Set 1",
					Name:        "Test Set 1",
					Description: "Description of the Set 1",
					Date:        time.Time{},
				},
				Elements: []*Element{
					{
						Metadata: Metadata{
							ID:          "element-1",
							Name:        "Element 1",
							Description: "Description of the element",
							Date:        time.Time{},
						},
						Values: []*ElementValue[any]{
							&ElementValue[any]{Value: "Test"},
							&ElementValue[any]{Value: 1},
						},
					},
				},
			},
		},
		Collections: []*Collection{
			{
				Metadata: Metadata{
					ID:          "collection-2",
					Name:        "Test Collection 2",
					Description: "Description of the collection 2",
					Date:        time.Time{},
				},
				Sets: []*Set{
					{
						Metadata: Metadata{
							ID:          "Set 2",
							Name:        "Test Set 2",
							Description: "Description of the set 2",
							Date:        time.Time{},
						},
						Elements: []*Element{
							{
								Metadata: Metadata{
									ID:          "element-2",
									Name:        "Element 2",
									Description: "Description of the element",
									Date:        time.Time{},
								},
								Values: []*ElementValue[any]{
									&ElementValue[any]{Value: 1},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.NotNil(t, collection)

}

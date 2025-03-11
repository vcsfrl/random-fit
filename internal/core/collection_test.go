package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	collection := Collection{
		Metadata: Metadata{
			ID:          "coll-pick-1",
			Name:        "Lotto number picks",
			Description: "Users monthly Lotto Number picks",
			Date:        time.Now(),
		},
		Collections: []*Collection{
			{
				Metadata: Metadata{
					ID:          "coll-pick-1-u1",
					Name:        "User Lotto Numbers",
					Description: "User Lotto Number picks",
					Date:        time.Now(),
				},
				Sets: []*Set{
					{
						Metadata: Metadata{
							ID:          "set-pick-u1-1",
							Name:        "6/49",
							Description: "User Lotto Number picks for 6/49",
							Date:        time.Now(),
						},
						Elements: []*Element{
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-1",
									Name:        "Numbers",
									Description: "6 numbers out of 49",
								},
								Values: []*ElementValue[any]{
									{Value: 1},
									{Value: 2},
									{Value: 3},
									{Value: 4},
									{Value: 5},
									{Value: 6},
								},
							},
							{
								Metadata: Metadata{
									ID:          "element-pick-u1-2",
									Name:        "Lucky Number",
									Description: "Lucky Number for 6/49 draw",
									Date:        time.Now(),
								},
								Values: []*ElementValue[any]{
									{Value: 25600},
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

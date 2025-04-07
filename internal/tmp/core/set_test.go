package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSet_String(t *testing.T) {
	set := Set{
		Metadata: Metadata{
			ID:          "Set 1",
			Name:        "Test",
			Description: "Description of the set",
			Date:        time.Now(),
		},
		Elements: []*Element[any]{
			{
				Metadata: Metadata{
					ID:          "element-1",
					Name:        "Element 1",
					Description: "Description of the element",
					Date:        time.Now(),
				},
				Values: []any{
					"Test", 1,
				},
			},
			{
				Metadata: Metadata{
					ID:          "element-2",
					Name:        "Element 2",
					Description: "Description of the element",
					Date:        time.Now(),
				},
				Values: []any{
					"Test", 2,
				},
			},
		},
	}

	assert.NotNil(t, set)
}

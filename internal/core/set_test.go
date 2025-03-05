package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSet_String(t *testing.T) {
	set := Set{
		Identity: Identity{
			ID:          "Set 1",
			Name:        "Test",
			Description: "Description of the set",
			Date:        time.Now(),
		},
		Elements: []Element{
			{
				Identity: Identity{
					ID:           "element-1",
					DefinitionID: "definition-1",
					Name:         "Element 1",
					Description:  "Description of the element",
					Date:         time.Now(),
				},
				Values: []fmt.Stringer{&ElementValue[string]{Value: "Test"}, &ElementValue[int]{Value: 1}},
			},
			{
				Identity: Identity{
					ID:           "element-2",
					DefinitionID: "definition-2",
					Name:         "Element 2",
					Description:  "Description of the element",
					Date:         time.Now(),
				},
				Values: []fmt.Stringer{&ElementValue[string]{Value: "Test"}, &ElementValue[int]{Value: 1}},
			},
		},
	}

	assert.NotNil(t, set)
}

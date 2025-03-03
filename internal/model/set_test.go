package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSet_String(t *testing.T) {
	set := Set{
		ID:          "Set 1",
		Name:        "Test",
		Description: "Description of the set",
		Elements: []*Element{
			{
				ID:     "element-1",
				Name:   "Element 1",
				Values: []fmt.Stringer{&ElementValue[string]{Value: "Test"}, &ElementValue[int]{Value: 1}},
				Date:   time.Now(),
			},
			{
				ID:     "element-2",
				Name:   "Element 2",
				Values: []fmt.Stringer{&ElementValue[string]{Value: "Test"}, &ElementValue[int]{Value: 1}},
				Date:   time.Now(),
			},
		},
	}

	assert.NotNil(t, set)
}

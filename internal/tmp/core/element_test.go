package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestElement_String(t *testing.T) {
	element := Element[int]{
		Metadata: Metadata{
			ID:          "element-1",
			Name:        "Element 1",
			Description: "Description of the element",
			Date:        time.Now(),
		},
		Values: []int{
			1, 2,
		},
	}

	assert.NotNil(t, element)
}

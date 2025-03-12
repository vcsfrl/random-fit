package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestElement_String(t *testing.T) {
	element := Element{
		Metadata: Metadata{
			ID:          "element-1",
			Name:        "Element 1",
			Description: "Description of the element",
			Date:        time.Now(),
		},
		Values: []*ElementValue[any]{
			&ElementValue[any]{Value: "Test"},
			&ElementValue[any]{Value: 12},
		},
	}

	assert.NotNil(t, element)
}

func TestElementValue_String(t *testing.T) {
	intValue := ElementValue[any]{Value: 1}
	assert.Equal(t, "1", intValue.String())

	stringValue := ElementValue[string]{Value: "Test"}
	assert.Equal(t, "Test", stringValue.String())
}

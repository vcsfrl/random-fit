package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElement_String(t *testing.T) {
	intElement := Element{
		ID:   "element-1",
		Name: "Test",
		Values: []fmt.Stringer{
			&ElementValue[string]{Value: "Test"},
			&ElementValue[int]{Value: 2},
		},
	}

	assert.Equal(t, "Test 2", intElement.String())
}

func TestElementValue_String(t *testing.T) {
	intValue := ElementValue[int]{Value: 1}
	assert.Equal(t, "1", intValue.String())

	stringValue := ElementValue[string]{Value: "Test"}
	assert.Equal(t, "Test", stringValue.String())
}

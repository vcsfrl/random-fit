package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElement_String(t *testing.T) {
	intElement := Element[int]{
		ID:    1,
		Name:  "Test",
		Value: ElementValue[int]{Value: 1},
	}

	assert.Equal(t, "1", intElement.String())
}

func TestElementValue_String(t *testing.T) {
	intValue := ElementValue[int]{Value: 1}
	assert.Equal(t, "1", intValue.String())

	stringValue := ElementValue[string]{Value: "Test"}
	assert.Equal(t, "Test", stringValue.String())
}

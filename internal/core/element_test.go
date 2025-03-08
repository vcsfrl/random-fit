package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestElement_String(t *testing.T) {
	intElement := Element{
		Metadata: Metadata{
			ID:           "element-1",
			DefinitionID: "definition-1",
			Name:         "Element 1",
			Description:  "Description of the element",
			Date:         time.Now(),
		},
		Values: []fmt.Stringer{
			&ElementValue[string]{Value: "Test"},
			&ElementValue[int]{Value: 12},
		},
		ValueSeparator: "-",
	}

	assert.Equal(t, "Test-12", intElement.String())
}

func TestElementValue_String(t *testing.T) {
	intValue := ElementValue[int]{Value: 1}
	assert.Equal(t, "1", intValue.String())

	stringValue := ElementValue[string]{Value: "Test"}
	assert.Equal(t, "Test", stringValue.String())
}

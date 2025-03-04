package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	collection := Collection{
		Identity: Identity{
			ID:           "collection-1",
			DefinitionID: "definition-1",
			Name:         "Test",
			Description:  "Description of the collection",
			Date:         time.Time{},
		},
		Sets: nil,
	}

	assert.NotNil(t, collection)

}

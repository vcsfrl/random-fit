package internal

import (
	"random-fit/internal/core"
	"time"
)

func Generate(definition core.ElementDefinition) core.Element {
	return core.Element{
		Identity: core.Identity{
			ID:           "element-1",
			DefinitionID: definition.Identity.ID,
			Name:         definition.Identity.Name,
			Description:  definition.Identity.Description,
			Date:         time.Time{},
		},
		Values: nil,

		ValueSeparator: " ",
	}
}

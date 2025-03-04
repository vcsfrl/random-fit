package internal

import (
	"random-fit/internal/model"
	"time"
)

func Generate(definition model.ElementDefinition) model.Element {
	return model.Element{
		Identity: model.Identity{
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

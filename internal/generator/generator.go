package generator

import (
	"random-fit/internal/core"
	"time"
)

type Generator struct {
}

func (g *Generator) Collection(definition core.CollectionDefinition) core.Collection {
	return core.Collection{}
}

func (g *Generator) Set(definition core.SetDefinition) core.Set {
	return core.Set{}
}

func (g *Generator) Element(definition core.ElementDefinition) core.Element {
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

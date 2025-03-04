package internal

import (
	"random-fit/internal/model"
	"time"
)

func Generate(definition model.ElementDefinition) model.Element {
	return model.Element{
		ID:             definition.ID,
		Name:           definition.Name,
		Values:         nil,
		Date:           time.Time{},
		ValueSeparator: " ",
	}
}

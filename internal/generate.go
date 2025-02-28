package internal

import "random-fit/internal/model"

func Generate[T any](definition model.ElementDefinition[T]) model.Element[T] {
	return model.Element[T]{
		ID:   0,
		Name: "",
	}
}

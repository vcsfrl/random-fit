package core

// Element is a generic type for an element.
type Element[T any] struct {
	Metadata Metadata
	Values   []T
}

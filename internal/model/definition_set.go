package model

// SetDefinition is a generic type for the definition of a set.
type SetDefinition struct {
	ID          string
	Name        string
	Description string
	Elements    []ElementDefinition
}

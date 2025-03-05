package model

// SetDefinition is a generic type for the definition of a set.
type SetDefinition struct {
	Identity DefinitionIdentity

	Elements []ElementDefinition
}

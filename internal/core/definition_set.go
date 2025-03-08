package core

// SetDefinition is a generic type for the definition of a set.
type SetDefinition struct {
	Metadata DefinitionMetadata

	Elements []ElementDefinition
}

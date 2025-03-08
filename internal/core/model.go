package core

import "time"

type Metadata struct {
	ID           string
	DefinitionID string
	Name         string
	Description  string
	Date         time.Time
}

type DefinitionMetadata struct {
	ID          string
	Name        string
	Description string
}

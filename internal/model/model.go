package model

import "time"

type Identity struct {
	ID           string
	DefinitionID string
	Name         string
	Description  string
	Date         time.Time
}

type DefinitionIdentity struct {
	ID          string
	Name        string
	Description string
}

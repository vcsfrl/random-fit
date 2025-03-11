package core

type Collection struct {
	Metadata Metadata

	Sets        []*Set
	Collections []*Collection
}

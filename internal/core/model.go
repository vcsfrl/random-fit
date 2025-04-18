package core

import "time"

type Metadata struct {
	ID       string
	ParentID string
	Details  string
	Date     time.Time
}

package marianatek

import (
	"encoding/json"
)

type ClassSessionTypeRelationships struct {
}

type ClassSessionTypeAttributes struct {
	Duration    int    `json:"duration"`
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ClassSessionType struct {
	Relationships ClassSessionTypeRelationships `json:"relationships"`
	Attributes    ClassSessionTypeAttributes    `json:"attributes"`
	Type          string                        `json:"type"`
	ID            int64                         `json:"id,string"`
}

func NewClassSessionType(e *Envelope) *ClassSessionType {
	s := &ClassSessionType{
		Type: e.Type,
		ID:   e.ID,
	}
	json.Unmarshal(e.Attributes, &s.Attributes)
	json.Unmarshal(e.Relationships, &s.Relationships)
	return s
}

package marianatek

import "encoding/json"

type LayoutRelationships struct {
	Classroom DataStruct     `json:"classroom"`
	Region    DataStruct     `json:"region"`
	Location  DataStruct     `json:"location"`
	Spots     DataListStruct `json:"spots"`
	Site      DataStruct     `json:"site"`
}

type LayoutAttributes struct {
	Listed   bool   `json:"listed"`
	Capacity int    `json:"capacity"`
	Name     string `json:"name"`
	Format   string `json:"format"`
}

type Layout struct {
	Relationships LayoutRelationships `json:"relationships"`
	Attributes    LayoutAttributes    `json:"attributes"`
	Type          string              `json:"type"`
	ID            int64               `json:"id"`
}

func NewLayout(e *Envelope) *Layout {
	s := &Layout{
		Type: e.Type,
		ID:   e.ID,
	}
	json.Unmarshal(e.Attributes, &s.Attributes)
	json.Unmarshal(e.Relationships, &s.Relationships)
	return s
}

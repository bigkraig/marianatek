package marianatek

import (
	"encoding/json"
)

type SpotRelationships struct {
	Classroom DataStruct `json:"classroom"`
	Layout    DataStruct `json:"layout"`
	Region    DataStruct `json:"region"`
	Site      DataStruct `json:"site"`
	Location  DataStruct `json:"location"`
	SpotType  DataStruct `json:"spot_type"`
}

type SpotAttributes struct {
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Enabled     bool   `json:"enabled"`
	YPosition   string `json:"y_position"`
	XPosition   string `json:"x_position"`
	Listed      bool   `json:"listed"`
	RowPosition int    `json:"row_position"`
}

type Spot struct {
	Relationships SpotRelationships `json:"relationships"`
	Attributes    SpotAttributes    `json:"attributes"`
	Type          string            `json:"type"`
	ID            int64             `json:"id,string"`
}

func NewSpot(e *Envelope) *Spot {
	s := &Spot{
		Type: e.Type,
		ID:   e.ID,
	}
	json.Unmarshal(e.Attributes, &s.Attributes)
	json.Unmarshal(e.Relationships, &s.Relationships)
	return s
}

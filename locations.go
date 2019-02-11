package marianatek

import (
	"context"
	"encoding/json"
	"fmt"
)

type LocationsService service

type LocationRelationships struct {
	Classrooms               DataListStruct `json:"classrooms"`
	AddonsProductCollection  DataStruct     `json:"addons_product_collection"`
	Region                   DataStruct     `json:"region"`
	Site                     DataStruct     `json:"site"`
	DefaultProductCollection DataStruct     `json:"default_product_collection"`
	Partner                  DataStruct     `json:"partner"`
}

type LocationAttributes struct {
	AddressLine2  string `json:"address_line2"`
	StateProvince string `json:"state_province"`
	AddressLine1  string `json:"address_line1"`
	Name          string `json:"name"`
	City          string `json:"city"`
	EmailAddress  string `json:"email_address"`
	Longitude     string `json:"longitude"`
	PostalCode    string `json:"postal_code"`
	Listed        bool   `json:"listed"`
	Latitude      string `json:"latitude"`
	Timezone      string `json:"timezone"`
	PhoneNumber   string `json:"phone_number"`
	Description   string `json:"description"`
}

type Location struct {
	Relationships LocationRelationships `json:"relationships"`
	Attributes    LocationAttributes    `json:"attributes"`
	Type          string                `json:"type"`
	ID            int64                 `json:"id,string"`
}

func NewLocation(e *Envelope) *Location {
	s := &Location{
		Type: e.Type,
		ID:   e.ID,
	}
	json.Unmarshal(e.Attributes, &s.Attributes)
	json.Unmarshal(e.Relationships, &s.Relationships)
	return s
}

func (s *LocationsService) List(ctx context.Context) ([]*Location, *Includes, error) {
	u := fmt.Sprintf("locations")

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var locations []*Location
	resp, err := s.client.Do(ctx, req, &locations)
	if err != nil {
		return nil, nil, err
	}

	return locations, resp.Includes, nil
}

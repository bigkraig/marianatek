package marianatek

import "encoding/json"

type ProductCollectionRelationships struct {
	ProductCollectionToProductComparisonAssignements DataListStruct `json:"product_collection_to_product_comparison_assignments"`
}

type ProductCollectionAttributes struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProductCollection struct {
	Relationships ProductCollectionRelationships `json:"relationships"`
	Attributes    ProductCollectionAttributes    `json:"attributes"`
	Type          string                         `json:"type"`
	ID            int64                          `json:"id,string"`
}

func NewProductCollection(e *Envelope) *ProductCollection {
	s := &ProductCollection{
		Type: e.Type,
		ID:   e.ID,
	}
	json.Unmarshal(e.Attributes, &s.Attributes)
	json.Unmarshal(e.Relationships, &s.Relationships)
	return s
}

package marianatek

import "encoding/json"

type Envelope struct {
	Type          string `json:"type"`
	ID            int64  `json:"id,string"`
	Attributes    json.RawMessage
	Relationships json.RawMessage
}

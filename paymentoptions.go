package marianatek

import (
	"context"
	"fmt"
)

type PaymentOptionsService service

type PaymentOptionAttributes struct {
	Name               string      `json:"name"`
	ExpirationDatetime interface{} `json:"expiration_datetime"`
	Count              int64       `json:"count"`
	PaymentOptionID    int64       `json:"payment_option_id"`
	PaymentOptionType  string      `json:"payment_option_type"`
	PaymentDescription string      `json:"payment_description"`
	IsActive           bool        `json:"is_active"`
	ErrorMessage       string      `json:"error_message"`
	ErrorType          string      `json:"error_type"`
}

type PaymentOptionRelationships struct {
	PaymentOption DataStruct `json:"payment_option"`
}

type PaymentOption struct {
	Attributes    PaymentOptionAttributes    `json:"attributes"`
	Relationships PaymentOptionRelationships `json:"relationships"`
	Type          string                     `json:"type"`
	ID            string                     `json:"id"`
}

type PaymentOptionsListOptions struct {
	ClassSession     int64 `url:"class_session,omitempty"`
	PageSize         int   `url:"page_size,omitempty"`
	ReservedForGuest bool  `url:"reserved_for_guest,omitempty"`
	User             int64 `url:"user,omitempty"`
}

func (s *PaymentOptionsService) Get(ctx context.Context, opt PaymentOptionsListOptions) ([]*PaymentOption, *Includes, error) {
	u := fmt.Sprintf("payment_options")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var paymentOptions []*PaymentOption
	resp, err := s.client.Do(ctx, req, &paymentOptions)
	if err != nil {
		return nil, nil, err
	}

	return paymentOptions, resp.Includes, nil
}

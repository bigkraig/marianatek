package marianatek

import (
	"context"
)

type ReservationsService service

type ReservationRelationships struct {
	ClassSession DataStruct `json:"class_session"`
	Spot         DataStruct `json:"spot"`
	User         DataStruct `json:"user"`
}

type ReservationAttributes struct {
	BookedBy         string                     `json:"booked_by"`
	BookedFor        string                     `json:"booked_for"`
	Status           string                     `json:"status"`
	CreationDate     interface{}                `json:"creation_date"`
	CancelDate       interface{}                `json:"cancel_date"`
	ReservationType  string                     `json:"reservation_type"`
	CheckInDate      interface{}                `json:"check_in_date"`
	ReservedForGuest bool                       `json:"reserved_for_guest"`
	WaitlistWeight   interface{}                `json:"waitlist_weight"`
	GuestEmail       string                     `json:"guest_email"`
	PaymentOptions   []ReservationPaymentOption `json:"payment_options"`
}

type ReservationPaymentOption struct {
	Count           int64  `json:"count"`
	PaymentOptionID int64  `json:"payment_option_id"`
	Type            string `json:"type"`
}

type Reservation struct {
	Relationships ReservationRelationships `json:"relationships"`
	Attributes    ReservationAttributes    `json:"attributes"`
	Type          string                   `json:"type,omitempty"`
	ID            int64                    `json:"id,string,omitempty"`
}

type ReserveOptions struct {
	Type               string
	PaymentOptionCount int64
	PaymentOptionID    int64
	PaymentOptionType  string
	ClassSessionID     int64
	SpotID             int64
	UserID             int64
}

func (s *ReservationsService) Reserve(ctx context.Context, opts *ReserveOptions) (*Reservation, *Includes, error) {
	data := struct {
		Data Reservation `json:"data"`
	}{
		Reservation{
			Type: "reservations",
			Attributes: ReservationAttributes{
				ReservationType: opts.Type,
				PaymentOptions: []ReservationPaymentOption{
					ReservationPaymentOption{
						Count:           opts.PaymentOptionCount,
						PaymentOptionID: opts.PaymentOptionID,
						Type:            opts.PaymentOptionType,
					},
				},
			},
			Relationships: ReservationRelationships{
				ClassSession: DataStruct{
					Data{
						Type: "class_sessions",
						ID:   opts.ClassSessionID,
					},
				},
				Spot: DataStruct{
					Data{
						Type: "spots",
						ID:   opts.SpotID,
					},
				},
				User: DataStruct{
					Data{
						Type: "users",
						ID:   opts.UserID,
					},
				},
			},
		},
	}

	req, err := s.client.NewRequest("POST", "reservations", data)
	if err != nil {
		return nil, nil, err
	}

	rResp := new(Reservation)
	resp, err := s.client.Do(ctx, req, rResp)
	if err != nil {
		return nil, nil, err
	}

	return rResp, resp.Includes, nil
}

package marianatek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// CancelPenalty represents the response from the cancel penalty endpoint
type CancelPenalty struct {
	IsPenaltyCancel bool    `json:"is_penalty_cancel"`
	Message         *string `json:"message"`
}

// GetCancelPenalty checks if cancelling a reservation will incur a penalty
func (s *ReservationsService) GetCancelPenalty(ctx context.Context, reservationID int64) (*CancelPenalty, error) {
	u := fmt.Sprintf("customer/v1/me/reservations/%d/cancel_penalty", reservationID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var penalty CancelPenalty
	if err := json.Unmarshal(body, &penalty); err != nil {
		return nil, err
	}

	return &penalty, nil
}

// Cancel cancels a reservation by ID
func (s *ReservationsService) Cancel(ctx context.Context, reservationID int64) error {
	u := fmt.Sprintf("customer/v1/me/reservations/%d/cancel", reservationID)

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return err
	}

	return nil
}

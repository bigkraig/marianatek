package marianatek

import (
	"context"
	"fmt"
	"time"
)

type ClassSessionsService service

type ClassSessionRelationship struct {
	AvailablePrimarySpots   DataListStruct `json:"available_primary_spots"`
	AvailableSecondarySpots DataListStruct `json:"available_secondary_spots"`
	BookingWindows          DataListStruct `json:"booking_windows"`
	CancelledByUser         DataStruct     `json:"cancelled_by_user"`
	Classroom               DataStruct     `json:"classroom"`
	ClassSessionNotes       DataListStruct `json:"class_session_notes"`
	ClassSessionType        DataStruct     `json:"class_session_type"`
	EmployeePublicProfiles  DataListStruct `json:"employee_public_profiles"`
	Instructors             DataListStruct `json:"instructors"`
	Layout                  DataStruct     `json:"layout"`
	Location                Data           `json:"location"`
	Reservations            DataListStruct `json:"reservations"`
	SpotHolds               DataListStruct `json:"spot_holds"`
	SubstitutedInstructors  DataListStruct `json:"substituted_instructors"`
	Tags                    DataListStruct `json:"tags"`
}

type ClassSessionAttributes struct {
	StartDatetime                time.Time `json:"start_datetime"`
	Capacity                     int       `json:"capacity"`
	StandardReservationUserCount int       `json:"standard_reservation_user_count"`
	CheckedInUserCount           int       `json:"checked_in_user_count"`
	StandbyAvailability          int       `json:"standby_availability"`
	StartTime                    string    `json:"start_time"`
	RecurringStatus              string    `json:"recurring_status"`
	WaitlistReservationUserCount int       `json:"waitlist_reservation_user_count"`
	StandbyCapacity              int       `json:"standby_capacity"`
	Public                       bool      `json:"public"`
	SlotPrice                    []struct {
		Count    int `json:"count"`
		SlotType int `json:"slot_type"`
	} `json:"slot_price"`
	InBookingWindow      bool        `json:"in_booking_window"`
	VipUserCount         int         `json:"vip_user_count"`
	EndDatetime          time.Time   `json:"end_datetime"`
	CancellationDatetime interface{} `json:"cancellation_datetime"`
	ArchivedAt           interface{} `json:"archived_at"`
	AvailableSpots       []int64     `json:"available_spots"`
	PublicNote           interface{} `json:"public_note"`
	StartDate            string      `json:"start_date"`
	FirstTimeUserCount   int         `json:"first_time_user_count"`
}

type ClassSession struct {
	Relationships ClassSessionRelationship `json:"relationships"`
	Attributes    ClassSessionAttributes   `json:"attributes"`
	Type          string                   `json:"type"`
	ID            int64                    `json:"id,string"`
}

type ClassSessionListOptions struct {
	Include  string `url:"include,omitempty"`
	Location int64  `url:"location,omitempty"`
	MaxDate  string `url:"max_date,omitempty"`
	MinDate  string `url:"min_date,omitempty"`
	PageSize string `url:"page_size,omitempty"`
}

func (s *ClassSessionsService) List(ctx context.Context, opt ClassSessionListOptions) ([]*ClassSession, *Includes, error) {
	u := fmt.Sprintf("class_sessions")
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var classSessions []*ClassSession
	resp, err := s.client.Do(ctx, req, &classSessions)
	if err != nil {
		return nil, nil, err
	}

	return classSessions, resp.Includes, nil
}

func (s *ClassSessionsService) Get(ctx context.Context, class int64) (*ClassSession, *Includes, error) {
	u := fmt.Sprintf("class_sessions/%v?page_size=100&include=available_primary_spots%%2Cavailable_secondary_spots%%2Cstandby_availability%%2Clayout%%2Clocation.addons_product_collection", class)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var classSession *ClassSession
	resp, err := s.client.Do(ctx, req, &classSession)
	if err != nil {
		return nil, nil, err
	}

	return classSession, resp.Includes, nil
}

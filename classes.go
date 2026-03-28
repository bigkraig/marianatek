package marianatek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type ClassesService service

type ClassType struct {
	ID                string `json:"id"`
	Description       string `json:"description"`
	Duration          int    `json:"duration"`
	DurationFormatted string `json:"duration_formatted"`
	IsLiveStream      bool   `json:"is_live_stream"`
	Name              string `json:"name"`
}

type Classroom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type InstructorPhoto struct {
	LargeURL     *string `json:"large_url"`
	ThumbnailURL string  `json:"thumbnail_url"`
}

type Instructor struct {
	ID             string          `json:"id"`
	Bio            string          `json:"bio"`
	InstagramHandle string         `json:"instagram_handle"`
	InstagramURL   string          `json:"instagram_url"`
	Name           string          `json:"name"`
	PhotoURLs      InstructorPhoto `json:"photo_urls"`
	SpotifyURL     *string         `json:"spotify_url"`
}

type ClassTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClassLocation struct {
	ID                                string   `json:"id"`
	AddressLineOne                    string   `json:"address_line_one"`
	AddressLineTwo                    string   `json:"address_line_two"`
	AddressLineThree                  string   `json:"address_line_three"`
	City                              string   `json:"city"`
	CurrencyCode                      string   `json:"currency_code"`
	Description                       string   `json:"description"`
	Email                             string   `json:"email"`
	GateGeoCheckInByDistance          bool     `json:"gate_geo_check_in_by_distance"`
	GeoCheckInDistance                int      `json:"geo_check_in_distance"`
	IsNewsletterSubscriptionPreChecked bool    `json:"is_newsletter_subscription_pre_checked"`
	Latitude                          string   `json:"latitude"`
	Longitude                         string   `json:"longitude"`
	PrimaryLanguage                   string   `json:"primary_language"`
	Name                              string   `json:"name"`
	PaymentGatewayType                string   `json:"payment_gateway_type"`
	PhoneNumber                       string   `json:"phone_number"`
	PostalCode                        string   `json:"postal_code"`
	StateProvince                     string   `json:"state_province"`
	Timezone                          string   `json:"timezone"`
	FormattedAddress                  []string `json:"formatted_address"`
	Listed                            bool     `json:"listed"`
	Region                            Region   `json:"region"`
}

type SpotOptions struct {
	PrimaryAvailability   int `json:"primary_availability"`
	PrimaryCapacity       int `json:"primary_capacity"`
	SecondaryAvailability int `json:"secondary_availability"`
	SecondaryCapacity     int `json:"secondary_capacity"`
	StandbyAvailability   int `json:"standby_availability"`
	StandbyCapacity       int `json:"standby_capacity"`
	WaitlistAvailability  int `json:"waitlist_availability"`
	WaitlistCapacity      int `json:"waitlist_capacity"`
}

type Class struct {
	ID                       string         `json:"id"`
	AvailableSpotCount       int            `json:"available_spot_count"`
	BookingStartDatetime     time.Time      `json:"booking_start_datetime"`
	IsLateCancelOverridden   bool           `json:"is_late_cancel_overridden"`
	Capacity                 int            `json:"capacity"`
	ClassTags                []ClassTag     `json:"class_tags"`
	ClassType                ClassType      `json:"class_type"`
	Classroom                Classroom      `json:"classroom"`
	ClassroomName            string         `json:"classroom_name"`
	GeoCheckInEndDatetime    time.Time      `json:"geo_check_in_end_datetime"`
	GeoCheckInStartDatetime  time.Time      `json:"geo_check_in_start_datetime"`
	Instructors              []Instructor   `json:"instructors"`
	InLiveStreamWindow       bool           `json:"in_live_stream_window"`
	IsCancelled              bool           `json:"is_cancelled"`
	IsRemainingSpotCountPublic bool         `json:"is_remaining_spot_count_public"`
	IsUserGuestReserved      bool           `json:"is_user_guest_reserved"`
	IsUserReserved           bool           `json:"is_user_reserved"`
	IsUserWaitlisted         bool           `json:"is_user_waitlisted"`
	IsFreeClass              bool           `json:"is_free_class"`
	LayoutFormat             string         `json:"layout_format"`
	Location                 ClassLocation  `json:"location"`
	Name                     string         `json:"name"`
	ShouldShowStudioAddress  bool           `json:"should_show_studio_address"`
	SpotOptions              SpotOptions    `json:"spot_options"`
	StartDate                string         `json:"start_date"`
	StartTime                string         `json:"start_time"`
	StartDatetime            time.Time      `json:"start_datetime"`
	Status                   *string        `json:"status"`
	WaitlistCount            *int           `json:"waitlist_count"`
	Reservations             []interface{}  `json:"reservations"`
}

type ClassListOptions struct {
	MinStartDate string `url:"min_start_date,omitempty"`
	MaxStartDate string `url:"max_start_date,omitempty"`
	PageSize     int    `url:"page_size,omitempty"`
	Location     int64  `url:"location,omitempty"`
	Region       int64  `url:"region,omitempty"`
}

type ClassesResponse struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []*Class `json:"results"`
	Meta     struct {
		Pagination struct {
			Count   int `json:"count"`
			Pages   int `json:"pages"`
			Page    int `json:"page"`
			PerPage int `json:"per_page"`
		} `json:"pagination"`
	} `json:"meta"`
}

func (s *ClassesService) List(ctx context.Context, opt ClassListOptions) ([]*Class, error) {
	u := "customer/v1/classes"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, err
	}

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

	var classesResp ClassesResponse
	if err := json.Unmarshal(body, &classesResp); err != nil {
		return nil, err
	}

	return classesResp.Results, nil
}

func (s *ClassesService) Get(ctx context.Context, classID string) (*Class, *Includes, error) {
	u := fmt.Sprintf("customer/v1/classes/%s", classID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return nil, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var class Class
	if err := json.Unmarshal(body, &class); err != nil {
		return nil, nil, err
	}

	return &class, nil, nil
}

// ClassSpot represents an available spot for a class
type ClassSpot struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SpotType string `json:"spot_type"`
}

// ClassSpotsResponse represents the response from the spots endpoint
type ClassSpotsResponse struct {
	Count    int          `json:"count"`
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Results  []*ClassSpot `json:"results"`
}

// GetSpots retrieves available spots for a class
func (s *ClassesService) GetSpots(ctx context.Context, classID string) ([]*ClassSpot, error) {
	u := fmt.Sprintf("customer/v1/classes/%s", classID)

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

	var spotsResp ClassSpotsResponse
	if err := json.Unmarshal(body, &spotsResp); err != nil {
		return nil, err
	}

	return spotsResp.Results, nil
}

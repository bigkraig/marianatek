package marianatek

import (
	"context"
	"fmt"
	"time"
)

type UsersService service

type UserRelationships struct {
	HomeLocation  DataStruct     `json:"home_location"`
	ProfileValues DataListStruct `json:"profile_values"`
	LastRegion    DataStruct     `json:"last_region"`
	Tags          DataListStruct `json:"tags"`
}

type UserAttributes struct {
	LastName                     string        `json:"last_name"`
	EmergencyContactPhone        string        `json:"emergency_contact_phone"`
	ThirdPartySync               bool          `json:"third_party_sync"`
	ApplyAccountBalanceToFees    bool          `json:"apply_account_balance_to_fees"`
	IsStaff                      bool          `json:"is_staff"`
	PostalCode                   string        `json:"postal_code"`
	FullName                     string        `json:"full_name"`
	WaiverSignedDatetime         time.Time     `json:"waiver_signed_datetime"`
	AddressLine2                 string        `json:"address_line2"`
	City                         string        `json:"city"`
	FirstName                    string        `json:"first_name"`
	AddressLine1                 string        `json:"address_line1"`
	SignedWaiver                 bool          `json:"signed_waiver"`
	EmergencyContactRelationship interface{}   `json:"emergency_contact_relationship"`
	PhoneNumber                  string        `json:"phone_number"`
	HasVipTagCache               bool          `json:"has_vip_tag_cache"`
	AccountBalance               float64       `json:"account_balance"`
	IsMinimal                    bool          `json:"is_minimal"`
	Permissions                  []interface{} `json:"permissions"`
	StateProvince                string        `json:"state_province"`
	Gender                       interface{}   `json:"gender"`
	Email                        string        `json:"email"`
	EmergencyContactEmail        interface{}   `json:"emergency_contact_email"`
	EmergencyContactName         string        `json:"emergency_contact_name"`
	Country                      interface{}   `json:"country"`
	BirthDate                    interface{}   `json:"birth_date"`
	MarketingOptIn               bool          `json:"marketing_opt_in"`
}

type User struct {
	Attributes    UserAttributes    `json:"attributes"`
	Relationships UserRelationships `json:"relationships"`
	Type          string            `json:"type"`
	ID            int64             `json:"id,string"`
}

func (s *UsersService) Get(ctx context.Context, user string) (*User, *Includes, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v", user)
	} else {
		u = "self"
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(User)
	resp, err := s.client.Do(ctx, req, uResp)
	if err != nil {
		return nil, nil, err
	}

	return uResp, resp.Includes, nil
}

package marianatek

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClassSesssionsService_Get(t *testing.T) {
	client, mux, _, response, teardown := setup("class_sessions.ID.GET.json")
	defer teardown()

	id := int64(27257)

	mux.HandleFunc(fmt.Sprintf("/class_sessions/%v", id), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, response)
	})

	session, includes, err := client.ClassSessions.Get(context.Background(), id)
	if err != nil {
		t.Errorf("ClassSessions.Get returned error: %v", err)
	}

	want := int64(244346)
	if !reflect.DeepEqual(session.ID, want) {
		t.Errorf("ClassSessions.Get returned %+v, want %+v", session.ID, want)
	}

	want = int64(27257)
	var spot *Spot
	for _, spot = range includes.Spots {
		if spot.Attributes.Name == "SpotName" {
			break
		}
	}

	if !reflect.DeepEqual(spot.ID, want) {
		t.Errorf("ClassSessions.Get Spot was %+v, want %+v", spot.ID, want)
	}
}

func TestClassSesssionsService_List(t *testing.T) {
	client, mux, _, response, teardown := setup("class_sessions.GET.json")
	defer teardown()

	mux.HandleFunc("/class_sessions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, response)
	})

	opt := ClassSessionListOptions{
		Include:  "available_primary_spots%2Cavailable_secondary_spots%2Cclass_session_type%2Cstandby_availability%2Cin_booking_window",
		Location: 11111,
		MaxDate:  "date",
		MinDate:  "date",
	}

	sessions, _, err := client.ClassSessions.List(context.Background(), opt)
	if err != nil {
		t.Errorf("ClassSessions.List returned error: %v", err)
	}

	want := 10
	if !reflect.DeepEqual(len(sessions), want) {
		t.Errorf("ClassSessions.List returned %+v, want %+v", len(sessions), want)
	}
}

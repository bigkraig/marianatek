package marianatek

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestReservationsService_Reserve(t *testing.T) {
	client, mux, _, response, teardown := setup("reservations.POST.json")
	defer teardown()

	mux.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, response)
	})

	opts := &ReserveOptions{
		Type:               "standard",
		PaymentOptionCount: 1,
		PaymentOptionID:    244639,
		PaymentOptionType:  "membership",
		ClassSessionID:     int64(240443),
		SpotID:             int64(27258),
		UserID:             int64(18841),
	}

	reserved, _, err := client.Reservations.Reserve(context.Background(), opts)
	if err != nil {
		t.Errorf("Reservations.Reserve returned error: %v", err)
	}

	want := "standard"
	if !reflect.DeepEqual(reserved.Attributes.ReservationType, want) {
		t.Errorf("Reservations.Reserve returned %+v, want %+v", reserved.Type, want)
	}
}

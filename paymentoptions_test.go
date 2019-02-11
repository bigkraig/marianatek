package marianatek

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPaymentOptionsService_Get(t *testing.T) {
	client, mux, _, response, teardown := setup("payment_options.GET.json")
	defer teardown()

	mux.HandleFunc("/payment_options", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, response)
	})

	opts := PaymentOptionsListOptions{
		ClassSession:     int64(12345),
		PageSize:         100,
		ReservedForGuest: false,
		User:             int64(54321),
	}

	options, _, err := client.PaymentOptions.Get(context.Background(), opts)
	if err != nil {
		t.Errorf("PaymentOptions.Get returned error: %v", err)
	}

	want := 2
	if !reflect.DeepEqual(len(options), want) {
		t.Errorf("PaymentOptions.Get returned %+v, want %+v", len(options), want)
	}
}

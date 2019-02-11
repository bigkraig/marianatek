package marianatek

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestLocationsService_List(t *testing.T) {
	client, mux, _, response, teardown := setup("locations.GET.json")
	defer teardown()

	mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, response)
	})

	locations, _, err := client.Locations.List(context.Background())
	if err != nil {
		t.Errorf("Locations.List returned error: %v", err)
	}

	want := 2
	if !reflect.DeepEqual(len(locations), want) {
		t.Errorf("Locations.List returned %+v, want %+v", len(locations), want)
	}
}

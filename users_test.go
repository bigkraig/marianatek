package marianatek

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUsersService_Get(t *testing.T) {
	client, mux, _, response, teardown := setup("users.GET.json")
	defer teardown()

	mux.HandleFunc("/users/self", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, response)
	})

	user, _, err := client.Users.Get(context.Background(), "self")
	if err != nil {
		t.Errorf("Users.Get returned error: %v", err)
	}

	want := int64(11111)
	if !reflect.DeepEqual(user.ID, want) {
		t.Errorf("Users.Get returned %+v, want %+v", user.ID, want)
	}
}

package marianatek

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	baseURLPath = "/api"
)

func TestClientError(t *testing.T) {
	client, mux, _, response, teardown := setup("client.error.json")
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, response, 422)
	})

	req, _ := client.NewRequest("GET", "", nil)

	_, err := client.Do(context.Background(), req, nil)
	if err == nil {
		t.Errorf("client.Do did not return error")
	}

	want := 1
	if !reflect.DeepEqual(len(err.(*ErrorResponse).Errors["spot"]), want) {
		t.Errorf("client.Do returned %+v errors, wanted %+v: %v", len(err.(*ErrorResponse).Errors["spot"]), want, err.Error())
	}
}

func TestClientIncludes(t *testing.T) {
	client, mux, _, response, teardown := setup("client.included.json")
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, response)
	})

	req, _ := client.NewRequest("GET", "", nil)

	resp, err := client.Do(context.Background(), req, nil)
	if err != nil {
		t.Errorf("client.Do returned error: %v", err)
	}

	want := 1
	if !reflect.DeepEqual(len(resp.Includes.Spots), want) {
		t.Errorf("client.Do returned %+v Spot includes, want %+v", len(resp.Includes.Spots), want)
	}
}

func setup(jsonFile string) (client *Client, mux *http.ServeMux, serverURL, responseData string, teardown func()) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client = NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	if jsonFile != "" {
		b, err := ioutil.ReadFile(filepath.Join("tests", jsonFile))
		if err != nil {
			log.Fatal("unexpected error ", err, " during os.Open on: ", jsonFile)
		} else {
			responseData = string(b)
		}
	}

	return client, mux, server.URL, responseData, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

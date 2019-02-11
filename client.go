package marianatek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://barrysbootcamp.marianatek.com/api/"
	userAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1 Safari/605.1.15"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	client *http.Client
	common service

	ClassSessions  *ClassSessionsService
	Locations      *LocationsService
	PaymentOptions *PaymentOptionsService
	Reservations   *ReservationsService
	Users          *UsersService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.common.client = c

	c.ClassSessions = (*ClassSessionsService)(&c.common)
	c.Locations = (*LocationsService)(&c.common)
	c.PaymentOptions = (*PaymentOptionsService)(&c.common)
	c.Reservations = (*ReservationsService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			return nil, e
		}

		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	response := newResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.Unmarshal(response.Data, v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

type ErrorResponse struct {
	Response *http.Response      // HTTP response that caused this error
	Errors   map[string][]string `json:"errors"` // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors)
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

type apiResponse struct {
	Meta struct {
		Pagination struct {
			Count   int `json:"count"`
			PerPage int `json:"per_page"`
			Pages   int `json:"pages"`
			Page    int `json:"page"`
		} `json:"pagination"`
	} `json:"meta"`
	Data     json.RawMessage `json:"data"`
	Included json.RawMessage `json:"included"`
	// Links interface{} `json:"links"`
}

type Data struct {
	Type string `json:"type"`
	ID   int64  `json:"id,string"`
}

type DataStruct struct {
	Data Data `json:"data"`
}
type DataListStruct struct {
	Data []Data `json:"data"`
}

type Response struct {
	*http.Response

	Count   int
	PerPage int
	Pages   int
	Page    int

	Data     json.RawMessage
	Includes *Includes
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.populatePageValues()
	return response
}

func (r *Response) populatePageValues() {
	v := &apiResponse{}
	err := json.NewDecoder(r.Body).Decode(v)
	if err == io.EOF {
		err = nil // ignore EOF errors caused by empty response body
	}

	if err != nil {
		log.Fatal(err)
	}

	r.Count = v.Meta.Pagination.Count
	r.PerPage = v.Meta.Pagination.PerPage
	r.Pages = v.Meta.Pagination.Pages
	r.Page = v.Meta.Pagination.Page
	r.Data = v.Data

	r.populateIncluded(v.Included)
}

func (r *Response) populateIncluded(data json.RawMessage) {
	r.Includes = &Includes{}

	var envelopes []*Envelope
	if err := json.Unmarshal(data, &envelopes); err != nil {
		return
	}

	for _, env := range envelopes {
		switch env.Type {
		case "spots":
			r.Includes.Spots = append(r.Includes.Spots, NewSpot(env))
		case "locations":
			r.Includes.Locations = append(r.Includes.Locations, NewLocation(env))
		case "layouts":
			r.Includes.Layouts = append(r.Includes.Layouts, NewLayout(env))
		case "product_collections":
			r.Includes.ProductCollections = append(r.Includes.ProductCollections, NewProductCollection(env))
		case "class_session_types":
			r.Includes.ClassSessionTypes = append(r.Includes.ClassSessionTypes, NewClassSessionType(env))
		default:
			log.Printf("unknown message type: %q", env.Type)
		}
	}
}

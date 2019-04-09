package revgeo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Query struct has properties and methods to allow the query
// to be assembled prior to the reverse geocode query
type Query struct {
	lat     float64
	lng     float64
	zoom    int
	address int
	email   string
}

// Response struct contains result of a unmarshaled successful
// reverse geocode query
type Response struct {
	PlaceID     int    `json:"place_id"`
	Licence     string `json:"licence"`
	OsmType     string `json:"osm_type"`
	OsmID       int    `json:"osm_id"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	DisplayName string `json:"display_name"`
	Address     `json:"address"`
}

// Address struct embedded in Response struct
type Address struct {
	HouseNumber   string `json:"house_number"`
	Road          string `json:"road"`
	Suburb        string `json:"suburb"`
	Village       string `json:"village"`
	Town          string `json:"town"`
	City          string `json:"city"`
	County        string `json:"county"`
	StateDistrict string `json:"state_district"`
	State         string `json:"state"`
	Postcode      string `json:"postcode"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
}

// NewDecoder is a constructor it returns a Query instance
func NewDecoder() *Query {
	var q *Query
	// set zoom at 18, this gives maximum detail in the response
	//  which we can unmarshal and use as needed
	q = &Query{
		zoom: 18,
	}
	return q
}

// SetLatLng allows the lat and lng coordinates we want to reverse
// geocode to be configured
func (q *Query) SetLatLng(lat float64, lng float64) {
	if lat != 0 && lng != 0 {
		q.lat = lat
		q.lng = lng
	}
}

// SetZoom allows the zoom be edited default is max 18
func (q *Query) SetZoom(zoom int) {
	if zoom >= 0 && zoom <= 18 {
		q.zoom = zoom
	}
}

// IncludeAddress requests that a full line by line address is
// provided as part of the response
func (q *Query) IncludeAddress(include bool) {
	if include {
		q.address = 1
	}
}

// SetEmail allows us to pass an email address, useful as
// there are no API keys or authentication required to use the
// service. If there are problems, related to request volume
// email address will be used to make contact
func (q *Query) SetEmail(email string) {
	if email != "" {
		q.email = email
	}
}

// Decode takes the assembled Query and makes the request for
// the coordinates to be decoded
func (q *Query) Decode() (*Response, error) {
	var err error
	var r *Response
	query := "https://nominatim.openstreetmap.org/reverse?format=json&lat=%v&lon=%v&zoom=%v&addressdetails=%v"

	// client timeout 5 seconds
	ht := &http.Client{
		Timeout: time.Second * 5,
	}

	// if no lat or lng return with error
	if q.lat == 0 || q.lng == 0 {
		return r, errors.New("Decode: Cannot make query, Lat or Lng not set")
	}

	// build request
	req := fmt.Sprintf(query, q.lat, q.lng, q.zoom, q.address)

	// add email if present
	if q.email != "" {
		req = req + "&email=" + q.email
	}

	// make request
	resp, err := ht.Get(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	// error if not 200 OK
	if resp.StatusCode != 200 {
		return r, errors.New("Decode: Response not 200 OK")
	}

	// process body and unmarshal
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	r, err = q.unmarshal(body, r)

	return r, nil
}

// helper for unmarshalling JSON, allows testing
func (q *Query) unmarshal(body []byte, r *Response) (*Response, error) {
	err := json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

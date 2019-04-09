package revgeo

import (
	"encoding/json"
	"testing"
)

// set this to true to use the HTTPClient for testing and compare
// JSON response to known response
// increases test coverage from cica 40% to >80%
const testHTTPClient = false

var testResponse = []byte(`{
	"place_id": 95126103,
	"licence": "Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright",
	"osm_type": "way",
	"osm_id": 90394420,
	"lat": "52.54877605",
	"lon": "-1.81627033283164",
	"display_name": "137, Pilkington Avenue, Sutton Coldfield, Birmingham, West Midlands Combined Authority, West Midlands, England, B72 1LH, United Kingdom",
	"address": {
	"house_number": "137",
	"road": "Pilkington Avenue",
	"town": "Sutton Coldfield",
	"city": "Birmingham",
	"county": "West Midlands Combined Authority",
	"state_district": "West Midlands",
	"state": "England",
	"postcode": "B72 1LH",
	"country": "United Kingdom",
	"country_code": "gb"
	},
	"boundingbox": [
	"52.5487321",
	"52.5488299",
	"-1.8163514",
	"-1.8161885"
	]
}
`)

var remarshaled = []byte(`{"place_id":95126103,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"way","osm_id":90394420,"lat":"52.54877605","lng":"","display_name":"137, Pilkington Avenue, Sutton Coldfield, Birmingham, West Midlands Combined Authority, West Midlands, England, B72 1LH, UK","address":{"house_number":"137","road":"Pilkington Avenue","town":"Sutton Coldfield","city":"Birmingham","county":"West Midlands Combined Authority","state_district":"West Midlands","state":"England","postcode":"B72 1LH","country":"UK","country_code":"gb"}}`)

func (q *Query) GetLatLng() (float64, float64) {
	return q.lat, q.lng
}

func TestSetLatLng(t *testing.T) {
	lat := 52.54877605
	lng := -1.81627033283164
	d := NewDecoder()
	d.SetLatLng(lat, lng)

	if d.lat != lat {
		t.Errorf("Expected lat of %v, got %v", lat, d.lat)
	}

	if d.lat != lat {
		t.Errorf("Expected lng of %v, got %v", lng, d.lng)
	}
}

func TestIncludeAddress(t *testing.T) {
	d := NewDecoder()
	d.IncludeAddress(true)
	if d.address != 1 {
		t.Errorf("Expected adddress equal to 1, got %v", d.address)
	}
}

func TestSetEmail(t *testing.T) {
	email := "hello@test.com"
	d := NewDecoder()
	d.SetEmail(email)
	if d.email != email {
		t.Errorf("Expected email to be %s, got %s", email, d.email)
	}
}

func TestUnmarshal(t *testing.T) {
	placeID := 95126103
	city := "Birmingham"
	postcode := "B72 1LH"

	d := NewDecoder()
	r := new(Response)

	r, err := d.unmarshal(testResponse, r)
	if err != nil {
		t.Error(err)
	}

	if r.PlaceID != placeID {
		t.Errorf("Expected placeID to be %v, got %v", placeID, r.PlaceID)
	}

	if r.Address.City != city {
		t.Errorf("Expected city to be %s, got %s", city, r.City)
	}

	if r.Address.Postcode != postcode {
		t.Errorf("Expected postcode to be %s, got %s", postcode, r.Postcode)
	}
}

func TestDecode(t *testing.T) {
	if testHTTPClient {
		lat := 52.54877605
		lng := -1.81627033283164
		d := NewDecoder()
		d.SetLatLng(lat, lng)
		d.IncludeAddress(true)
		r, err := d.Decode()
		if err != nil {
			t.Error(err)
		}

		data, err := json.Marshal(r)
		if err != nil {
			t.Error(err)
		}
		if string(data) != string(remarshaled) {
			t.Error("Mismatched JSON, expected to be the same")
		}
	}
}

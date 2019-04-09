## Reverse geocoding using Nominatim (openstreetmap.org)

Simple client library for the Nominatim service which provides free reverse geocoding of lat/lng coordinates to address data

## Usage

```Go
rg := revgeo.NewDecoder()

// coordinates to query
rg.SetLatLng(52.54877605, -1.81627033283164)

// include detailed address in response
rg.IncludeAddress(true)

// include address, used if there are rate limiting issues
rg.SetEmail("me@example.com")

r, err := rg.Decode()
```

## Tests

Tests are basic. You can increase coverage by enabling the HTTP client in ```internal_test.go```, which will make live queries to the Nominatim service as part of the test run. Should be used prudently.

```Go
const testHTTPClient = false // set to true
```

## License

This package is MIT licensed. Nominatim has its own [licensing and Acceptable Use Policy (AUP)](https://operations.osmfoundation.org/policies/nominatim/). 
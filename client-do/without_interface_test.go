package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type testingRoundtripper struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (t testingRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTrip(req)
}

// Test_WithoutInterface brings way more meaning to our tests since we're able to intercept the response
// and construct http responses to validate the handling request logic
func Test_WithoutInterface(t *testing.T) {
	tcs := []struct {
		url         string
		expectedErr error
		expected    ipInfoMeta
		roundTrip   func(req *http.Request) (*http.Response, error)
	}{
		// Happy path
		{
			url:         "https://ipinfo.io",
			expectedErr: nil,
			expected: ipInfoMeta{
				IP:           "14.102.144.0",
				City:         "Kuala Lumpur",
				Region:       "Pahang",
				Country:      "MY",
				Location:     "2.5000,112.5000",
				Organization: "AS45352 IP ServerOne Solutions Sdn Bhd",
			},
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{
						"ip": "14.102.144.0",
						"city": "Kuala Lumpur",
						"region": "Pahang",
						"country": "MY",
						"loc": "2.5000,112.5000",
						"org": "AS45352 IP ServerOne Solutions Sdn Bhd"
					  }`))),
				}, nil
			},
		},
		// Service unavailable
		{
			url:         "https://ipinfo.io",
			expectedErr: httpError{503, ""},
			expected:    ipInfoMeta{},
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 503,
					Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{
						"error": "Service Unavailable"
					  }`))),
				}, nil
			},
		},
		// {
		// 	// Got XML back for some reason
		// 	roundTrip: func(req *http.Request) (*http.Response, error) {
		// 		return &http.Response{
		// 			StatusCode: 200,
		// 			Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`
		// 			<?xml version="1.0" encoding="UTF-8" ?>
		// 			<root>
		// 				<ip>14.102.144.0</ip>
		// 				<city>Kuala Lumpur</city>
		// 				<region>Pahang</region>
		// 				<country>MY</country>
		// 				<loc>2.5000,112.5000</loc>
		// 				<org>AS45352 IP ServerOne Solutions Sdn Bhd</org>
		// 			</root>`))),
		// 		}, nil
		// 	},
		// 	expectedErr: &json.SyntaxError{},
		// 	expected:    ipInfoMeta{},
		// },
	}

	for _, tc := range tcs {
		s := NewSender(&http.Client{
			Transport: testingRoundtripper{
				roundTrip: tc.roundTrip,
			},
		})

		// I will leave the URL empty here as
		// we construct the http response above
		ipMeta, err := handleRequestWithoutInterface(s, "")
		if err != tc.expectedErr {
			t.Errorf("expected error was %v but got %v", tc.expectedErr, err)
		}

		if !reflect.DeepEqual(ipMeta, tc.expected) {
			t.Errorf("expected ipMeta to equal:\n %v\n but got:\n %v", tc.expected, ipMeta)
		}
	}
}

package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

// StubSender implements the proxy interface so we can "test"
type StubSender struct {
	client *http.Client
}

// do won't make an HTTP request here but will simply return dummy data
// solely for the purpose of testing
func (s StubSender) do(URL string) ([]byte, error) {
	if URL == "https://ipinfo.io" {
		// return somewhere in Malaysia
		return []byte(`{
			"ip": "14.102.144.0",
			"city": "Kuala Lumpur",
			"region": "Pahang",
			"country": "MY",
			"loc": "2.5000,112.5000",
			"org": "AS45352 IP ServerOne Solutions Sdn Bhd"
		  }`), nil
	}

	if URL == "https://www.google.com" {
		// return empty JSON here
		return []byte(`{}`), nil
	}

	if URL == "https://www.iknowhowtotestinterfaces.com" {
		// return an error here
		return nil, http.ErrHandlerTimeout
	}

	return nil, nil
}

func Test_WithInterface(t *testing.T) {
	s := StubSender{
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}

	tcs := []struct {
		URL         string
		expected    ipInfoMeta
		expectedErr error
	}{
		{
			URL:         "https://ipinfo.io",
			expectedErr: nil,
			expected: ipInfoMeta{
				IP:           "14.102.144.0",
				City:         "Kuala Lumpur",
				Region:       "Pahang",
				Country:      "MY",
				Location:     "2.5000,112.5000",
				Organization: "AS45352 IP ServerOne Solutions Sdn Bhd",
			},
		},
		{
			URL:         "https://www.google.com",
			expectedErr: nil,
			expected: ipInfoMeta{
				IP:           "",
				City:         "",
				Region:       "",
				Country:      "",
				Location:     "",
				Organization: "",
			},
		},
		{
			URL:         "https://www.iknowhowtotestinterfaces.com",
			expectedErr: http.ErrHandlerTimeout,
			expected:    ipInfoMeta{},
		},
	}

	for _, tc := range tcs {
		ipMeta, err := handleRequestWithInterface(s, tc.URL)

		if err != tc.expectedErr {
			t.Errorf("expected error was %v but got %v", tc.expectedErr, err)
		}

		if !reflect.DeepEqual(ipMeta, tc.expected) {
			t.Errorf("expected ipMeta to equal:\n %v\n but got:\n %v", tc.expected, ipMeta)
		}
	}
}

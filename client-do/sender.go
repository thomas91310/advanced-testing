package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpError struct {
	statusCode int
	url        string
}

func (h httpError) Error() string {
	return fmt.Sprintf("http error %v while getting %v", h.statusCode, h.url)
}

// Sender embeds an http client and its goal is to send http requests
// Sender implements the Proxyer interface
type Sender struct {
	client *http.Client
}

// NewSender creates a new sender
func NewSender(c *http.Client) Sender {
	return Sender{
		client: c,
	}
}

// do sends a GET request to URL
func (s Sender) do(URL string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// ignoring the body of the response when I don't get back a 200
	// don't copy this
	if resp.StatusCode != 200 {
		return nil, httpError{resp.StatusCode, URL}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

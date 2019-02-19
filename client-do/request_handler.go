package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type option string

const (
	withoutInterface option = "withoutInterface"
	withInterface    option = "withInterface"
)

func handleRequest(c *http.Client, URL string, opt option) (ipInfoMeta, error) {
	s := NewSender(c)

	switch opt {
	case withoutInterface:
		return handleRequestWithoutInterface(s, URL)
	case withInterface:
		return handleRequestWithInterface(s, URL)
	default:
		return ipInfoMeta{}, fmt.Errorf("unsupported option %v", opt)
	}
}

func handleRequestWithInterface(p proxyer, URL string) (ipInfoMeta, error) {
	b, err := p.do(URL)
	if err != nil {
		return ipInfoMeta{}, err
	}

	ipMeta := ipInfoMeta{}
	err = json.Unmarshal(b, &ipMeta)
	if err != nil {
		return ipInfoMeta{}, err
	}

	return ipMeta, nil
}

func handleRequestWithoutInterface(s Sender, URL string) (ipInfoMeta, error) {
	b, err := s.do(URL)
	if err != nil {
		return ipInfoMeta{}, err
	}

	ipMeta := ipInfoMeta{}
	err = json.Unmarshal(b, &ipMeta)
	if err != nil {
		return ipInfoMeta{}, err
	}

	return ipMeta, nil
}

package main

import "fmt"

// ipInfoMeta represents the JSON response I get back from ipinfo.io
type ipInfoMeta struct {
	IP           string `json:"ip"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Country      string `json:"country"`
	Location     string `json:"loc"`
	PostalCode   string `json:"postal"`
	Organization string `json:"org"`
}

// Print prints a ipInfoMeta object
func (i ipInfoMeta) print() {
	fmt.Println("ip:", i.IP)
	fmt.Println("city:", i.City)
	fmt.Println("region:", i.Region)
	fmt.Println("country:", i.Country)
	fmt.Println("loc:", i.Location)
	fmt.Println("postal:", i.PostalCode)
	fmt.Println("org:", i.Organization)
}

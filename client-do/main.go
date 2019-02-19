package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("I am a program that sends a request to ipinfo.io")
	fmt.Println("------------------------------------------------")

	c := &http.Client{
		Timeout: 1 * time.Second,
	}

	ipMeta, err := handleRequest(c, "https://ipinfo.io", withoutInterface)
	if err != nil {
		panic(err)
	}
	ipMeta.print()

	fmt.Println("------------------------------------------------")
	fmt.Println("done!")
}

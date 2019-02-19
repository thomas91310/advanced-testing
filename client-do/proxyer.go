package main

type proxyer interface {
	do(url string) ([]byte, error)
}

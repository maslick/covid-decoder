package main

import "github.com/maslick/covid-decoder/src"

func main() {
	server := src.RestController{Service: &src.Service{}}
	server.Start()
}

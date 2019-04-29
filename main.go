package main

import (
	"fmt"
	"net/http"
)

func main() {

	lb := ParseConfig()
	lb.initialize()

	http.HandleFunc("/", lb.handleRequest)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Print("Failed to start server")
	}

}

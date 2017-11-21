package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

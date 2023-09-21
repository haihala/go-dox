package main

import (
	"fmt"
	"log"
	"net/http"
)

func parrot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", parrot)
	mux.HandleFunc("/hello", hello)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

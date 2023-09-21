package main

import (
	"fmt"
	"log"
	"net/http"
)

func parrot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You sent '%s'", r.URL.Path[1:])
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", parrot)
	mux.HandleFunc("/hello", hello)

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

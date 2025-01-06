package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!!!!!!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe("0.0.0.0:8080", mux)
    if err != nil {
        log.Fatal(err)
    }
}
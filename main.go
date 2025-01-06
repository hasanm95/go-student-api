package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hasanm95/go-student-api/internal/config"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Student API!!!"))
}

func main() {
	// Load config from file
	cfg := config.MustLoad()

	fmt.Println(cfg)


	// Setup router
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(cfg.HTTPServer.Addr, mux)
    if err != nil {
        log.Fatal(err)
    }

	fmt.Println("Server started on port 8080")

}
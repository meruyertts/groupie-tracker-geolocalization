package main

import (
	"fmt"
	"net/http"

	"groupie-tracker-geolocalization/handlers"
)

func main() {
	mux := http.NewServeMux()
	styles := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", styles))
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/artist/", handlers.ArtistHandler)
	fmt.Println("Server launched at http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("This host is already run")
	}
}

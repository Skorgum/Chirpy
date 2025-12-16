package main

import (
	"fmt"
	"net/http"
)

func main() {
	//create new http.ServeMux
	mux := http.NewServeMux()

	//create new http.Server struct
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	//start the server
	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}

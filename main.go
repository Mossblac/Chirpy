package main

import (
	"fmt"
	"net/http"

	"github.com/Mossblac/Chirpy/ext"
)

func main() {
	mux := http.NewServeMux()

	server := ext.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := http.ListenAndServe(server.Addr, server.Handler)
	fmt.Printf("%v", err)
}

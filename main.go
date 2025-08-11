package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8080", mux)
	fmt.Printf("%v", err)

}

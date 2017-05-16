package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting up...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, cotinuous deployment.. and again")
	})
	http.ListenAndServe(":80", nil)
}

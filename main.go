package main

import (
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello Docker</h1>"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8000", nil)
}

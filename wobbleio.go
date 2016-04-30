package main

import (
	"net/http"
	"fmt"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello Docker</h1>"))
	fmt.Println("Listen/n")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8000", nil)
}

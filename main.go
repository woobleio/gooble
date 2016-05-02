package main

import (
	"os"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>y</h1>"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}

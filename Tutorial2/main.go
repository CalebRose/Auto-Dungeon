package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func about(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>I AM NOW A GOD. YOUR GOD. I DESERVE THIS.</h1>")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	fmt.Println("Server starting...")
	http.ListenAndServe(":3000", nil)
}
package main

import (
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello!")
	})


	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":80", nil))
}

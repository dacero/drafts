package main

import (
	"net/http"
	"text/template"
	"log"
)


func serveTemplate(w http.ResponseWriter, t *http.Request) {
	template, err := template.ParseFiles("./template.html")
	if err != nil {
		log.Printf("Error when opening the template: %s", err)
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Printf("Error when executing the template: %s", err)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", serveTemplate)
	log.Fatal(http.ListenAndServe(":80", nil))
}

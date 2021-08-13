package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func serveTemplate(w http.ResponseWriter, t *http.Request) {
	fmt.Println("Listing drafts directory")
	myDir, _ := listDraftDirectory("./drafts/")
	fmt.Printf("Now the items are: %s\n", len(myDir.Files))
	
	template, err := template.ParseFiles("./templates/template.gohtml")
	if err != nil {
		log.Printf("Error when opening the template: %s", err)
	}
	
	d, err := ioutil.ReadFile("./drafts/draft.md")
	if err != nil {
		log.Printf("Error when opening the draft: %s", err)
	}
	draft := Draft{Body: string(d)}
		
	err = template.Execute(w, draft)
	if err != nil {
		log.Printf("Error when executing the template: %s", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", serveTemplate)
	r.HandleFunc("/view", viewDraft)
	r.HandleFunc("/dir", viewDirectory)
	log.Fatal(http.ListenAndServe(":80", r))
}

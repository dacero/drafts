package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"github.com/gorilla/mux"
)

type Draft struct {
	Body    string
}

func (d Draft) HTMLBody() string {
	groomedString := strings.ReplaceAll(d.Body, "\r\n", "\n")
	unsafe := blackfriday.Run([]byte(groomedString))
	output := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(output)
}

func serveTemplate(w http.ResponseWriter, t *http.Request) {
	template, err := template.ParseFiles("./template.gohtml")
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

func viewDraft(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	filename := r.FormValue("filename")

	template, err := template.ParseFiles("./template.gohtml")
	if err != nil {
		fmt.Fprintf(w, "Error when opening the template: %s", err)
		return
	}
	
	d, err := ioutil.ReadFile("./drafts/" + filename)
	if err != nil {
		fmt.Fprintf(w, "Error when opening the draft: %s", err)
		return
	}
	draft := Draft{Body: string(d)}
		
	err = template.Execute(w, draft)
	if err != nil {
		fmt.Fprintf(w, "Error when executing the template: %s", err)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", serveTemplate)
	r.HandleFunc("/view", viewDraft)
	log.Fatal(http.ListenAndServe(":80", r))
}

package main

import (
	"net/http"
	"text/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/microcosm-cc/bluemonday"
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

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", serveTemplate)
	log.Fatal(http.ListenAndServe(":80", nil))
}

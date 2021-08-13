package main

import (
	"strings"
	"fmt"
	"text/template"
	"io/ioutil"
	"net/http"
	"log"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
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

func viewDraft(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("ParseForm() err: %v", err)
		return
	}
	filepath := r.FormValue("filepath")

	template, err := template.ParseFiles("./templates/template.gohtml")
	if err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("Error when opening the draft template: %s", err)
		return
	}
	
	d, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(w, "Error when opening the draft: %s", filepath)
		log.Printf("Error when opening the draft: %s", err)
		return
	}
	draft := Draft{Body: string(d)}
		
	err = template.Execute(w, draft)
	if err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("Error when executing the draft template: %s", err)
		return
	}
}


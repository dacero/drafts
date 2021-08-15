package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/sessions"
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

func viewDraftHandler(store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth, _ := checkAuthorization(w, r, store); !auth {
			http.Redirect(w, r, "/static/door.html", http.StatusFound)
		}
		
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
	})
}


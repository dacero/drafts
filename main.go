package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dir?dir=drafts", http.StatusFound)
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/view", viewDraft)
	r.HandleFunc("/dir", viewDirectory)
	log.Fatal(http.ListenAndServe(":80", r))
}

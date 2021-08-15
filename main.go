package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/door.html", http.StatusFound)
}

func main() {
	var store *sessions.CookieStore
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)
	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0, //60 * 15,
		HttpOnly: true,
	}

	
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/view", viewDraftHandler(store))
	r.HandleFunc("/dir", viewDirectory)
	r.HandleFunc("/authenticate", authenticate(store))
	log.Fatal(http.ListenAndServe(":80", r))
}

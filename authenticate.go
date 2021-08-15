package main

import (
	"net/http"
	"log"
	"os"
	"fmt"

	"github.com/gorilla/sessions"
)

func authenticate(store *sessions.CookieStore) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "drafts-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		secret := r.FormValue("secret")
		if secret != os.Getenv("DRAFTS_KEY") {
			session.Values["authenticated"] = false
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
		
		session.Values["authenticated"] = true
		err = session.Save(r, w)
		if err != nil {
			log.Print("Internal Server Error when saving the session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		http.Redirect(w, r, "/dir?dir=drafts", http.StatusFound)
	})
}

func checkAuthorization(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) (bool, error) {
	if store == nil {
		return true, nil
	}
	
	session, err := store.Get(r, "drafts-session")
	if err != nil {
		return false, err
	}
	auth := session.Values["authenticated"]
	if auth == nil {
		return false, nil
	}
	if !auth.(bool) {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

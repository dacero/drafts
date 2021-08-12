package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"github.com/gorilla/mux"
)

type Draft struct {
	Body    string
}

type DraftFile struct {
	Filename   string
	Path       string // always relative to basepath
	IsDir      bool
}

type DraftsDirectory struct {
	Name    string
	Files   []DraftFile
}

func (d Draft) HTMLBody() string {
	groomedString := strings.ReplaceAll(d.Body, "\r\n", "\n")
	unsafe := blackfriday.Run([]byte(groomedString))
	output := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(output)
}

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

func viewDirectory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	dirpath := r.FormValue("dir")
	
	draftsDir, err := listDraftDirectory("./" + dirpath)
	if err != nil {
		fmt.Fprintf(w, "Error when opening the directory: %s", err)
		return
	}

	template, err := template.ParseFiles("./templates/dir.gohtml")
	if err != nil {
		log.Printf("Error when opening the template: %s", err)
	}
	err = template.Execute(w, draftsDir)
	if err != nil {
		log.Printf("Error when executing the template: %s", err)
	}
}


func viewDraft(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	filepath := r.FormValue("filepath")

	template, err := template.ParseFiles("./templates/template.gohtml")
	if err != nil {
		fmt.Fprintf(w, "Error when opening the template: %s", err)
		return
	}
	
	d, err := ioutil.ReadFile(filepath)
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

func listDraftDirectory(dirpath string) (DraftsDirectory, error) {
	draftFiles := []DraftFile{}
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return DraftsDirectory{}, err
	}
	for _, f := range files {
		if (filepath.Ext(f.Name()) == ".md" || f.IsDir()) {
			draftFiles = append(draftFiles, DraftFile{Filename: f.Name(),
				Path: filepath.Join(dirpath, f.Name()),
				IsDir: f.IsDir()})
		}
	}
	d := DraftsDirectory{Name: "Drafts", Files: draftFiles}
	return d, nil
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", serveTemplate)
	r.HandleFunc("/view", viewDraft)
	r.HandleFunc("/dir", viewDirectory)
	log.Fatal(http.ListenAndServe(":80", r))
}

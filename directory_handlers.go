package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"time"
	"text/template"
	"log"
	"path/filepath"
	"sort"
)

type DraftFile struct {
	Filename   string
	Path       string // always relative to basepath
	IsDir      bool
	ModTime    time.Time
}

type DraftsDirectory struct {
	Name    string
	Files   []DraftFile
}

func viewDirectory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("ParseForm() err: %v", err)
		return
	}
	dirpath := r.FormValue("dir")
	
	draftsDir, err := listDraftDirectory("./" + dirpath)
	if err != nil {
		fmt.Fprintf(w, "Error when opening the directory: %s", dirpath)
		log.Printf("Error when opening the directory: %s", err)
		return
	}

	template, err := template.ParseFiles("./templates/dir.gohtml")
	if err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("Error when opening the directory template: %s", err)
		return
	}
	err = template.Execute(w, draftsDir)
	if err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("Error when executing the directory template: %s", err)
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
			filename := f.Name()
			extension := filepath.Ext(f.Name())
			draftFiles = append(draftFiles, DraftFile{Filename: filename[0:len(filename)-len(extension)],
				Path: filepath.Join(dirpath, f.Name()),
				IsDir: f.IsDir(),
				ModTime: f.ModTime()})
		}
	}
	sort.Slice(draftFiles, func(i, j int) bool {
		return draftFiles[i].ModTime.After(draftFiles[j].ModTime)
	})
	d := DraftsDirectory{Name: "Drafts", Files: draftFiles}
	return d, nil
}


package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func executeSyncScript() {
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		log.Printf("Could not find bash: %s", err)
	}
	
	cmd := &exec.Cmd {
		Path: bashPath,
		Args: []string{ bashPath, "sync.sh" },
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	if err = cmd.Run(); err != nil {
		log.Printf("Error running the sync script: %s", err)
	}
}


func syncDrafts(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Internal error. Please check log.")
		log.Printf("ParseForm() err: %v", err)
		return
	}
	dirpath := r.FormValue("dir")

	executeSyncScript()

	http.Redirect(w, r, "/dir?dir=" + dirpath, http.StatusFound)

	
}

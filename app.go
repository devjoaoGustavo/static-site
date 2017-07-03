package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")

	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Stat returns the FileInfo structure describing file. If there is an error, it will be of type *PathError
	info, err := os.Stat(fp)
	// 404 if the template does not exist
	if err != nil {
		/* IsNotFound returns a boolean indicating
		   whether the error is known to report
		   that a file or directory does not exist.

		func IsNotExist(err error) bool*/
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// func ParseFiles(filenames ...string) (*Template, error)
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

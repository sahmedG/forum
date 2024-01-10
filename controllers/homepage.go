package controllers

import (
	"html/template"
	"net/http"
)

// main page Handler
func RenderPage(w http.ResponseWriter, r *http.Request) {
	// Check if the request path is not the root path
	if r.URL.Path != "/" {
		HTTPErrorHandler(w, r, http.StatusNotFound)
		return
	}
	// Check if the request is not GET && NOT POST requests
	if r.Method != "GET" && r.Method != "POST" {
		HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"static/html/index.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		HTTPErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		HTTPErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

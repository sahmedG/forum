package controllers

import (
	"html/template"
	"net/http"
)

// Category page Handler
func RenderCreatePostPage(w http.ResponseWriter, r *http.Request) {

	// Accept method GET only
	if r.Method != http.MethodGet {
		HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	files := []string{
		"static/html/createpost.html",
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

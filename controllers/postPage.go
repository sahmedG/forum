package controllers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/pkgs/funcs"
)

// Category page Handler
func RenderPostPage(w http.ResponseWriter, r *http.Request) {
	// Check if the request is not GET && NOT POST requests
	if r.Method != "GET" {
		HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	postID := strings.TrimPrefix(r.URL.Path, "/post/")
	_, err := funcs.GetPostID(postID)
	if err != nil {
		http.Error(w, "This Post do not exist", http.StatusBadRequest)
		return
	}

	files := []string{
		"static/html/post.html", // need to create a post page
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

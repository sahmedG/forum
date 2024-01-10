package controllers

import (
	"database/sql"
	"html/template"
	"net/http"
)

// user page Handler (coming soon)
func RenderUserPage(w http.ResponseWriter, r *http.Request, data *sql.DB) {
	// Check if the request is not GET && NOT POST requests
	if r.Method != "GET" {
		HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	// userID := strings.TrimPrefix(r.URL.Path, "/user/")
	// _, err := funcs.GetCategoryID(userID)
	// if err != nil {
	// 	http.Error(w, "This Category do not exist", http.StatusBadRequest)
	// 	return
	// }
	/////////////////////////////////////////////////////////////////////////////////////////
	files := []string{
		"static/html/index.html", // need to creat a user page
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		HTTPErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		HTTPErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/pkgs/funcs"
	"net/http"
	"strings"
)

type posts_capsul_json struct {
	Posts []funcs.Post_json `json:"posts"`
}

// Category page Handler
func RenderCategoryPage(w http.ResponseWriter, r *http.Request, data *sql.DB) {
	// Check if the request is not GET && NOT POST requests
	if r.Method != "GET" {
		HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	categoryName := strings.TrimPrefix(r.URL.Path, "/category/")
	catID, err := funcs.GetCategoryID(categoryName)
	if err != nil {
		http.Error(w, "This Category do not exist", http.StatusBadRequest)
		return
	}
	postIDs, err := funcs.GetCategoryPosts(catID)
	if err != nil {
		http.Error(w, "This Category do not exist", http.StatusBadRequest)
		return
	}

	posts_arr, _ := funcs.Get_posts_of_category(postIDs)

	/* Used to encapsulate the struct into one struct that is used to construct JSON for sending to front-end */
	JSON_Response := posts_capsul_json{
		Posts: posts_arr,
	}
	/* Marshal the data into JSON format */
	jsonData, err := json.Marshal(JSON_Response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	/* Set the content type to JSON */
	w.Header().Set("Content-Type", "application/json")

	/* Write the JSON data to the response writer */
	w.Write(jsonData)
}

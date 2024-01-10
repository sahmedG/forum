package api

import (
	"encoding/json"
	"fmt"
	"forum/controllers"
	"forum/pkgs/funcs"
	"net/http"
	"strings"
)

// Category page Handler
func Categories_filter_handler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is not GET && NOT POST requests
	if r.Method != "GET" {
		controllers.HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	categoryName := strings.TrimPrefix(r.URL.Path, "/api/category/")
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

	userSession, valid := ValidateUser(w, r)
	posts_arr, _ := funcs.Get_posts_of_category(postIDs)

	if valid {
		for idx := range posts_arr {
			posts_arr[idx].IsLiked, _ = funcs.Post_is_liked_by_user(userSession.Get_UserID(), posts_arr[idx].Post_ID)
		}
	}

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

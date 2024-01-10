package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/controllers"
	"forum/pkgs/funcs"
)

// Category page Handler
func ByUser_filter_handler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is not GET && NOT POST requests
	if r.Method != http.MethodGet {
		controllers.HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	userSession, valid := ValidateUser(w, r)
	var posts_arr []funcs.Post_json
	if valid {
		onlyLiked := false
		if r.URL.Path == "/api/liked_by_user" {
			onlyLiked = true
		}
		postIDs, err := funcs.GetUserPosts(userSession.Get_UserID(), onlyLiked)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		posts_arr, _ = funcs.Get_posts_of_category(postIDs)
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

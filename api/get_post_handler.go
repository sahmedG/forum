package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"forum/controllers"
	"forum/pkgs/funcs"
)

// Category page Handler
func Get_post_handler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is not GET && NOT POST requests
	if r.Method != http.MethodGet {
		controllers.HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	postID := strings.TrimPrefix(r.URL.Path, "/api/post/")

	json_post, err := funcs.Get_Post(postID)
	json_post.Post_ID = postID // assign post id to the JSON

	if err != nil {
		http.Error(w, "This Post do not exist", http.StatusBadRequest)
		return
	}

	userSession, valid := ValidateUser(w, r)

	if valid {
		json_post.IsLiked, _ = funcs.Post_is_liked_by_user(userSession.Get_UserID(), postID)
	}
	/* Marshal the data into JSON format */
	jsonData, err := json.Marshal(json_post)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	/* Set the content type to JSON */
	w.Header().Set("Content-Type", "application/json")

	/* Write the JSON data to the response writer */
	w.Write(jsonData)
}

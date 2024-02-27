package api

import (
	"encoding/json"
	"fmt"
	"forum/pkgs/funcs"
	"io"
	"net/http"
	"strconv"
)

type Post_Comment struct {
	Post_id     int
	Start_index int
	End_index   int
}

type Comments_Container struct {
	Comments []funcs.CommentResults `json:"comments"`
}

func Serve_comments_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var data Post_Comment

	// Unmarshal the JSON data from the request body
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}
	userSession, valid := ValidateUser(w, r)

	comments, _ := funcs.GetComment(data.Post_id)

	if valid {
		for idx := range comments {
			comments[idx].IsLiked, _ = funcs.Comment_is_liked_by_user(userSession.Get_UserID(), strconv.Itoa(comments[idx].CommentID))
		}
	}

	comments_capsul := Comments_Container{
		Comments: comments,
	}

	/* Marshal the data into JSON format */
	jsonData, err := json.Marshal(comments_capsul)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	/* Set the content type to JSON */
	w.Header().Set("Content-Type", "application/json")

	/* Write the JSON data to the response writer */
	w.Write(jsonData)
}

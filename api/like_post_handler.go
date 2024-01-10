package api

import (
	"encoding/json"
	"fmt"
	"forum/pkgs/funcs"
	"io"
	"net/http"
)

type PostLike struct {
	PostID      int
	LikeDislike int
}

func LikesPostHandler(w http.ResponseWriter, r *http.Request) {
	// Handling only POST method requests
	if r.Method == http.MethodPost {
		userSession, valid := ValidateUser(w, r)
		if !valid {
			w.Write([]byte("Unauthorize access"))
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var data PostLike
		// Unmarshal the JSON data from the request body into 'data' variable of type Post
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
			return
		}

		err = funcs.PostLikes(userSession.Get_UserID(), data.PostID, data.LikeDislike)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
		w.Write([]byte("OK!"))

	} else {
		// Handle other HTTP methods or incorrect requests
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/pkgs/funcs"
	"io"
	"net/http"
)

type Comment struct {
	Post_id int
	Content string
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Handling only POST method requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	var data Comment

	// Unmarshal the JSON data from the request body into 'data' variable of type Post
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	// Remove leading and trailing white spaces from the Comment content and checks if it is empty and within the limits
	if err := CheckCommentData(data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	err = funcs.CreateComment(userSession.Get_UserID(), data.Post_id, data.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	w.Write([]byte("OK!"))
}

func CheckCommentData(data Comment) error {
	// Check for empty title and Comment content
	if checkEmpty(&data.Content) {
		return errors.New("Comment con not be empty")
	}

	// Check the length of the Comment content
	if len(data.Content) > 500 {
		return errors.New("Comment Content should be up to 500 characters long")
	}

	return nil
}

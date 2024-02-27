package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"forum/pkgs/funcs"
)

type Post struct {
	Post       string   `json:"Post"`
	Title      string   `json:"Title"`
	ImageFile  string   `json:"ImageFileBase64"`
	Categories []string `json:"Categories"`
}

func Create_Post(w http.ResponseWriter, r *http.Request) {
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
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	// 	return
	// }

	var data Post

	// Unmarshal the JSON data from the request body into 'data' variable of type Post
	// if err := json.Unmarshal(body, &data); err != nil {
	// 	fmt.Println("here1")
	// 	http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
	// 	return
	// }

	data.Post = r.FormValue("Post")
	data.Title = r.FormValue("Title")

	data.Categories = strings.Split(r.FormValue("Categories"),",") // Assuming Categories is an array
	// Access the image file

	file, imageFileHeader, err := r.FormFile("ImageFile")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Failed to get form file", http.StatusBadRequest)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// Remove leading and trailing white spaces from the title, post content and checks if it is empty and within the limits
	if err := CheckPostData(data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	categorysName := strings.TrimSpace(strings.Join(data.Categories, " "))
	if len(categorysName) < 1 {
		data.Categories = []string{"General"}
	}


	err = funcs.CreatePost(userSession.Get_UserID(), data.Title, data.Categories, data.Post, imageFileHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	// Process the image file if present
	if data.ImageFile != "" {
		// Decode the base64 image data and save it as needed
		decodedImage, err := base64.StdEncoding.DecodeString(data.ImageFile)
		if err != nil {
			http.Error(w, "Failed to decode base64 image data", http.StatusInternalServerError)
			return
		}
		fmt.Println(decodedImage)
		// Now you can use the 'decodedImage' as needed (e.g., save it to a file)
		// Example: ioutil.WriteFile("image.jpg", decodedImage, 0644)
	}

	fmt.Println("POST CREATED SUCCESS")

	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
}

func CheckPostData(data Post) error {
	// Check for empty title and post content
	if checkEmpty(&data.Post) || checkEmpty(&data.Title) {
		return errors.New("title and Post Content are required fields")
	}

	// Check the length of the post content
	if len(data.Post) > 10000 {
		return errors.New("post Content should be up to 10000 characters long")
	}

	// Check the length of the title
	if len(data.Title) > 100 {
		return errors.New("title should be up to 100 characters long")
	}

	return nil
}

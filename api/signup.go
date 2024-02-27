package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"forum/pkgs/funcs"
	"forum/pkgs/hashing"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	// Get method, serve the page
	if r.Method == http.MethodGet {
		// Parse the template
		tmpl, err := template.ParseFiles("static/html/signup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute the template with the data
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// Post method, serve the request
	} else if r.Method == http.MethodPost {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		var data SignUp_form
		// Unmarshal the JSON data from the request body
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
			return
		}

		if err := CheckSignUpData(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}
		// Hash the password before adding it
		hash, _ := hashing.HashPassword(data.Password) // ignore error for the sake of simplicity

		if err := funcs.AddUser(data.Uname, data.Email, hash, data.EX_ID); err != nil {
			w.WriteHeader(http.StatusConflict)
			io.WriteString(w, err.Error())
			return
		}
		successfulLogin := Successful_Login{
			User_id: UserSession.userID,
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successfulLogin)
		// io.WriteString(w, "Add user success")
		fmt.Printf("Name: %s, Email: %s, Password: %s", data.Uname, data.Email, data.Password)
	} else {
		// Handle other HTTP methods or incorrect requests
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Remove leading and trailing white spaces from the email and checks if it is empty
func checkEmpty(name *string) bool {
	*name = strings.TrimSpace(*name)
	return *name == ""
}

// checkPass checks if the given password string contains any whitespace characters.
// It returns true if there are whitespace characters, and false otherwise.
func checkPassWS(pass string) bool {
	return strings.ContainsAny(pass, " \t\n\r\v\f")
}

func CheckSignUpData(data *SignUp_form) error {
	// Remove leading and trailing white spaces from the email,user name and checks if it is empty
	if checkEmpty(&data.Email) || checkEmpty(&data.Uname) {
		return errors.New("username and email are required")
	}

	// checks if the password have any whitespace in it
	if checkPassWS(data.Password) {
		return errors.New("password cannot contain whitespaces")
	}

	// checks the length of the username
	if len(data.Uname) > 20 {
		return errors.New("username should be up to 20 characters long")
	}

	// checks the length of the email
	if len(data.Email) > 30 {
		return errors.New("email should be up to 30 characters long")
	}

	if len(data.Password) < 6 && len(data.Password) > 20 {
		return errors.New("password should be between 6 and 20 characters long")
	}
	return nil
}

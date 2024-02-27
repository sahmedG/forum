package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"

	"forum/pkgs/funcs"
	"forum/pkgs/hashing"
)
type Successful_Login struct {
	User_id int  `json:"user_id"`
	Success bool `json:"success"`
}
var UserSession Session
func LogIn(w http.ResponseWriter, r *http.Request) {
	// Get method, serve the page
	if r.Method == http.MethodGet {
		// Parse the template
		tmpl, err := template.ParseFiles("static/html/login.html")
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


		var data LogIn_form

		// Unmarshal the JSON data from the request body
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
			return
		}

		// Remove leading and trailing white spaces from the email and checks if it is empty
		if checkEmpty(&data.Email) || checkPassWS(data.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			successfulLogin := Successful_Login{
				User_id: -1,
				Success: false,
			}

			w.WriteHeader(http.StatusUnauthorized)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(successfulLogin)
			return
		}

		// get user id

		get_user_id, err := funcs.SelectUserID(data.Email, data.EX_ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			successfulLogin := Successful_Login{
				User_id: -1,
				Success: false,
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(successfulLogin)

			return
		}

		// // Check if the user already has an active session
		// if sessionToken, hasSession := UserHasSession(get_user_id); hasSession {
		// 	// Delete the old session to log out from the first browser
		// 	delete(Sessions, sessionToken)
		// }

		// Check if the user already has an active session
		if sessionToken, hasSession := UserHasSession(get_user_id); hasSession {
			// Log out the user from the first browser
			logoutErr := LogOutBySessionToken(w, sessionToken)
			if logoutErr != nil {
				http.Error(w, "Failed to log out from the first browser", http.StatusInternalServerError)
				return
			}
		}

		// Check if passwords matches
		hash_matched := hashing.CheckPasswordHash(data.Password, funcs.GetUserHash(get_user_id)) // ignore error for the sake of simplicity

		if !hash_matched {
			//	io.WriteString(w, "Pass doesn't match!")
			// send error data
			w.WriteHeader(http.StatusUnauthorized)
			successfulLogin := Successful_Login{
				User_id: -1,
				Success: false,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(successfulLogin)
			return
		}
		// Create a seesion for this user

		// Generate a new UUID
		userUUID, err := uuid.NewV4()
		if err != nil {
			// Handle the error
			fmt.Printf("error: %s\n", err)
		}

		// Associate the UUID with the user in your session or database
		UserSession = Session{
			userID:      get_user_id,
			SessionUUID: userUUID.String(),
			expiry:      time.Now().Add(3600 * time.Second),
		}
		Sessions[UserSession.SessionUUID] = UserSession

		// Set a cookie with a session token that can be used to authenticate access without logging in
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   UserSession.SessionUUID,
			Expires: UserSession.expiry,
		})

		fmt.Printf("UUID: %s\n", UserSession.SessionUUID)
		// send welcome data
		successfulLogin := Successful_Login{
			User_id: UserSession.userID,
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(successfulLogin)
		// A go routine to indicate that the session is expired
		go EXPIRED(UserSession)
	} else {
		// Handle other HTTP methods or incorrect requests
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func EXPIRED(userSession Session) {
	for !userSession.IsExpired() {
	}
	fmt.Printf("User %d token expired!\n", userSession.userID)
}

func IsLoggedIn(user string) bool {
	if _, ok := Sessions[user]; ok {
		if !Sessions[user].IsExpired() {
			fmt.Println("Already logged in")
			return true
		}
	}
	return false
}

func LogOutBySessionToken(w http.ResponseWriter, sessionToken string) error {
	// Get the session from the Sessions map
	if session, ok := Sessions[sessionToken]; ok {
		// Remove the session from the Sessions map
		delete(Sessions, sessionToken)
		// Expire the cookie
		cookie := &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
		}

		// Set the cookie in the HTTP response
		http.SetCookie(w, cookie)
		w.Header().Add("Set-Cookie", "session_token=; Max-Age=0; HttpOnly")
		fmt.Printf("User %d logged out!\n", session.userID)
	}

	return nil
}

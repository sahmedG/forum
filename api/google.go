package api

import (
	"bytes"
	"encoding/json"
	"forum/pkgs/funcs"
	"forum/pkgs/hashing"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	EX_ID    string `json:"id"`
	// Add other fields as needed
}

type JWTResponse struct {
	Token string `json:"credential"`
}

type User struct {
	Email string `json:"email"`
	ID    string `json:"sub"`
	// Add other user details as needed
}

const (
	googleClientID     = "653458676586-64o6vf69qvlhnluujbicgesa9fq5kb0f.apps.googleusercontent.com"
	googleClientSecret = "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"
	googleRedirectURI  = "https://localhost:443/google"
)

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	Code = r.URL.Query().Get("code")
	if Code == "" {
		authURL := "https://accounts.google.com/o/oauth2/auth"
		params := url.Values{}
		params.Add("client_id", googleClientID)
		params.Add("redirect_uri", googleRedirectURI)
		params.Add("response_type", "code")
		params.Add("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
		params.Add("access_type", "offline")

		w.Header().Set("Access-Control-Allow-Origin", "*")    // Allow all origins, you might want to change this to a specific origin
		w.Header().Set("Access-Control-Allow-Methods", "GET") // Allow only GET requests for simplicity

		http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusSeeOther)
	} else {
		HandleGoogleCallback(w, r)
	}
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	
	code := r.URL.Query().Get("code")

	tokenURL := "https://accounts.google.com/o/oauth2/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", googleRedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Fatal("Error getting token:", err)
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		log.Fatal("Error parsing token response:", err)
		http.Error(w, "Error parsing token response", http.StatusInternalServerError)
		return
	}

	accessToken := tokenData["access_token"].(string)
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error getting user info:", err)
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userInfoBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading user info response:", err)
		http.Error(w, "Error reading user info response", http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(userInfoBody, &userInfo); err != nil {
		log.Fatal("Error parsing user info:", err)
		http.Error(w, "Error parsing user info", http.StatusInternalServerError)
		return
	}
	hash, _ := hashing.HashPassword(userInfo["id"].(string))
	funcs.AddUser(userInfo["name"].(string), userInfo["email"].(string), hash, userInfo["id"].(string))

	userDetails := map[string]string{
		"email":    userInfo["email"].(string), // Replace with actual user details
		"password": userInfo["id"].(string),    // Replace with actual generated password
	}

	// Marshal user details to JSON
	userDetailsJSON, err := json.Marshal(userDetails)
	if err != nil {
		http.Error(w, "Failed to encode user details", http.StatusInternalServerError)
		return
	}
	// Create a new request to the login endpoint
	req, err = http.NewRequest("POST", "https://localhost:443/login", bytes.NewBuffer(userDetailsJSON))
	if err != nil {
		log.Println("Error creating request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json") // Set content type header

	// Send the request
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   UserSession.SessionUUID,
		Path:    "https://localhost:443",
		Expires: UserSession.expiry,
	})
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != http.StatusOK {
		// Handle non-200 status code
		http.Error(w, "Login failed", resp.StatusCode)
		return
	}
	http.Redirect(w, req, "https://localhost:443/", http.StatusSeeOther)
}

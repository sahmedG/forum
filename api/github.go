package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/pkgs/funcs"
	"forum/pkgs/hashing"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	githubClientID     = "1e7aa20b6eaf1123d2ab"
	githubClientSecret = "8111b1760e13ceddf2a3244e4618f92233bf7981"
	githubRedirectURI  = "https://localhost:443/github/callback"
)

func HandleGithubLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	Code = r.URL.Query().Get("code")
	if Code == "" {
		// authURL := "https://github.com/login/oauth/authorize"
		params := url.Values{}
		params.Add("client_id", githubClientID)
		params.Add("redirect_uri", githubRedirectURI)
		params.Add("scope", "user:email") // Include user:email scope
		params.Add("state", "github")
		// redirectURL := authURL + "?" + params.Encode()
		redirectURL := "https://github.com/login/oauth/authorize?client_id=1e7aa20b6eaf1123d2ab"
		Testing(w, r, redirectURL)
	} else {
		HandleGithubCallback(w, r)
	}
}

func Testing(w http.ResponseWriter, r *http.Request, redirect string) {
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func HandleGithubCallback(w http.ResponseWriter, r *http.Request) {

	fmt.Println(Code)

	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{}
	data.Set("code", Code)
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("redirect_uri", githubRedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Fatal("Error getting token:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return
	}

	log.Println("Token response:", string(body)) // Log token response for debugging

	// Extract access token from token response
	accessToken := ExtractAccessToken(string(body))
	if accessToken == "" {
		log.Fatal("Access token not found in token response")
		return
	}

	// Use the access token for further requests
	userInfoURL := "https://api.github.com/user"
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error getting user info:", err)
		return
	}
	defer resp.Body.Close()

	userInfoBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading user info response:", err)
		return
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(userInfoBody, &userInfo); err != nil {
		log.Fatal("Error parsing user info:", err)
		return
	}
	var id string
	switch v := userInfo["id"].(type) {
	case string:
		id = v
	case float64:
		id = fmt.Sprintf("%.0f", v)
	default:
		log.Fatal("Unable to convert 'id' to string")
		return
	}
	hash, _ := hashing.HashPassword(id)
	login_email := ""
	if userInfo["name"] == nil && userInfo["email"] == nil {
		funcs.AddUser(userInfo["login"].(string), userInfo["login"].(string), hash, "github")
		login_email = userInfo["login"].(string)
	} else if userInfo["name"] == nil && userInfo["email"] != nil {
		funcs.AddUser(userInfo["login"].(string), userInfo["email"].(string), hash, "github")
	} else if userInfo["name"] != nil && userInfo["email"] == nil {
		funcs.AddUser(userInfo["name"].(string), userInfo["login"].(string), hash, "github")
		login_email = userInfo["login"].(string)
	} else {
		funcs.AddUser(userInfo["name"].(string), userInfo["email"].(string), hash, "github")
		login_email = userInfo["email"].(string)
	}
	jsonResponse := map[string]string{
		"email":    login_email,
		"password": id,
	}

	// Marshal JSON payload
	jsonData, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Println("Error marshaling JSON data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create POST request to /login with JSON payload
	req, err = http.NewRequest("POST", "https://localhost:443/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json") // Set content type header

	// Make the request
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

func ExtractAccessToken(body string) string {
	params, err := url.ParseQuery(body)
	if err != nil {
		log.Println("Error parsing token response:", err)
		return ""
	}
	return params.Get("access_token")
}

// successfulLogin := Successful_Login{
// 	User_id: UserSession.userID,
// 	Success: true,
// }
// body, err := io.ReadAll(r.Body)
// if err != nil {
// 	http.Error(w, "Failed to read request body", http.StatusInternalServerError)
// 	return
// }
// var code CallbackGH
// if err := json.Unmarshal(body, &code); err != nil {
// 	http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
// 	return
// }

// code := r.URL.Query().Get("code")
// fmt.Println(code)
// w.Header().Set("Content-Type", "application/json")
// json.NewEncoder(w).Encode(successfulLogin)
// Details := LogIn_form{
// 	Email  : userInfo["email"].(string),
// 	Password: hash,
// 	EX_ID  :  "github",
// }
// jsonData, err := json.Marshal(Details)

// if err != nil {
// 	// Handle error if marshaling fails
// 	http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
// 	return
// }
// w.Header().Set("Content-Type", "application/json")

// // Write the JSON data to the response body
// _, err = w.Write(jsonData)
// if err != nil {
// 	http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
// 	return
// }

// LogIn(w, r)
// http.Redirect(w,r,"/",http.StatusSeeOther)

package api
import (
	"fmt"
	"net/http"
	"time"
)
func LogOut(w http.ResponseWriter, r *http.Request) {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No session token, user is not logged in
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Get the session from the Sessions map
	sessionToken := cookie.Value
	if session, ok := Sessions[sessionToken]; ok {
		// Remove the session from the Sessions map
		delete(Sessions, sessionToken)
		// Expire the cookie
		cookie.Expires = time.Now().Add(-time.Hour)
		http.SetCookie(w, cookie)
		fmt.Printf("User %d logged out!\n", session.userID)
	}
	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
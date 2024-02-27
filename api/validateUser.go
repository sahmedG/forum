package api

import (
	"net/http"
)

/* Takes cookies from clients and returns session,
* and boolean value indicating that user is logged or not */
func ValidateUser(w http.ResponseWriter, r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return nil, false
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}
	// Get cookie value
	session_token := cookie.Value
	// We then get the session from our session map
	userSession, exists := Sessions[session_token]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.IsExpired() {
		delete(Sessions, session_token)
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	// If the session is valid, return a welcome message to the user
	// fmt.Fprintf(w, "Welcome %d!", userSession.Get_UserID())
	return &userSession, true // return the session
}

// checkusersession checks if a user session is valid based on the provided userID.
// It returns an error if the session is not found or has expired, and nil otherwise.
func UserHasSession(userID int) (string, bool) {
	for tokin, session := range Sessions {
		if session.userID == userID {
			return tokin, true
		}
	}
	return "", false
}

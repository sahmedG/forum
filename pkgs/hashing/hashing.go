package hashing

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// // A middleware used for authenticate access to participate in the forum
// // All create, edit, post, comment handlers should be passed in this middleware handler
// func session_midddle(in_http http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			if err == http.ErrNoCookie {
// 				// If the cookie is not set, return an unauthorized status
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}
// 			// For any other type of error, return a bad request status
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		// Get cookie value
// 		session_token := cookie.Value

// 		// We then get the session from our session map
// 		userSession, exists := api.Sessions[session_token]
// 		if !exists {
// 			// If the session token is not present in session map, return an unauthorized error
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		// If the session is present, but has expired, we can delete the session, and return
// 		// an unauthorized status
// 		if userSession.IsExpired() {
// 			delete(api.Sessions, session_token)
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		// If the session is valid, return a welcome message to the user
// 		fmt.Fprintf(w, "Welcome %d!", userSession.Get_UserID())

// 		in_http.ServeHTTP(w, r)
// 	})
// }

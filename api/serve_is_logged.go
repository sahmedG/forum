package api

import "net/http"

func Serve_is_logged(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  _,islogged := ValidateUser(w, r)
  
  if islogged {
    w.Write([]byte("1"))
  } else {
    w.Write([]byte("0"))
  }
}

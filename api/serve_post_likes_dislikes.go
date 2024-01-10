package api

import (
	"forum/controllers"
	"forum/pkgs/funcs"
	"net/http"
  "strconv"
  "encoding/json"
  "fmt"
)

type like_dislike_capsul struct {
  Interactions funcs.LikeCounts `json:"interactions"`
}

func Serve_post_likes_dislikes (w http.ResponseWriter, r *http.Request) {

  if r.Method != http.MethodGet {
    controllers.HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)

    return
  }

 // Get the value of the "post_id" header
	postIDStr := r.Header.Get("post_id")

	// Convert the string to an integer
	postID, _ := strconv.Atoi(postIDStr)
  likes_dislikes_count,_ := funcs.CountPostLikes(postID)

  // For debug
  //fmt.Println(likes_dislikes_count)

/* Used to encapsulate the struct into one struct that is used to construct JSON for sending to front-end */
	JSON_Response := like_dislike_capsul{
		Interactions: likes_dislikes_count,
	}
	/* Marshal the data into JSON format */
	jsonData, err := json.Marshal(JSON_Response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	/* Set the content type to JSON */
	w.Header().Set("Content-Type", "application/json")

	/* Write the JSON data to the response writer */
	w.Write(jsonData)

}

func Serve_comm_likes_dislikes (w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
	  controllers.HTTPErrorHandler(w, r, http.StatusMethodNotAllowed)
  
	  return
	}
  
   // Get the value of the "comm_id" header
	  commIDStr := r.Header.Get("comm_id")
  
	  // Convert the string to an integer
	  commID, _ := strconv.Atoi(commIDStr)
	likes_dislikes_count,_ := funcs.CountCommentLikes(commID)
  
	// For debug
	//fmt.Println(likes_dislikes_count)
  
  /* Used to encapsulate the struct into one struct that is used to construct JSON for sending to front-end */
	  JSON_Response := like_dislike_capsul{
		  Interactions: likes_dislikes_count,
	  }
	  /* Marshal the data into JSON format */
	  jsonData, err := json.Marshal(JSON_Response)
	  if err != nil {
		  fmt.Println("Error marshaling JSON:", err)
		  return
	  }
  
	  /* Set the content type to JSON */
	  w.Header().Set("Content-Type", "application/json")
  
	  /* Write the JSON data to the response writer */
	  w.Write(jsonData)
  
  }
  
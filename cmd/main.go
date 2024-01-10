package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/api"
	"forum/controllers"
	"forum/pkgs/funcs"
)

func main() {
	funcs.Init()

	// Create a file server to serve static files (CSS, JS, images, etc.)
	fs := http.FileServer(http.Dir("static"))

	// Handle requests for files in the "/static/" path
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	/********************* API endpoints ************************/
	http.HandleFunc("/signup", api.SignUp) // Handle signup
	http.HandleFunc("/login", api.LogIn)   // Handle login
	http.HandleFunc("/logout", api.LogOut) //Handle logout

	http.HandleFunc("/api/islogged", api.Serve_is_logged)              // Check if user logged in or not
	http.HandleFunc("/api/created_by_user", api.ByUser_filter_handler) // posts filtering ex: /api/userposts (when the user is loggedin)
	http.HandleFunc("/api/liked_by_user", api.ByUser_filter_handler)

	http.HandleFunc("/api/create_post", api.Create_Post)     // create post
	http.HandleFunc("/api/posts", api.GetPostsHandler)       // Retrive posts as JSON
	http.HandleFunc("/api/post/", api.Get_post_handler)      // Retrive one post ex: /post/2
	http.HandleFunc("/api/likes_post", api.LikesPostHandler) // Handle Likes & Dislikes for Posts
	http.HandleFunc("/api/postlikes", api.Serve_post_likes_dislikes)

	http.HandleFunc("/api/add_comment", api.AddCommentHandler)     // Create comment
	http.HandleFunc("/api/comments", api.Serve_comments_handler)   // Serves post comments
	http.HandleFunc("/api/likes_comment", api.LikesCommentHandler) // Handle Likes & Dislikes for Posts
	http.HandleFunc("/api/commlikes", api.Serve_comm_likes_dislikes)

	http.HandleFunc("/api/create_category", api.Create_Category_Handler) // create category
	http.HandleFunc("/api/categories", api.Serve_categories_handler)     // Serves categories
	http.HandleFunc("/api/category/", api.Categories_filter_handler)     // Category filtering ex: /api/category/Technology

	/********************* END ************************/

	// Render pages
	http.HandleFunc("/", controllers.RenderPage)
	http.HandleFunc("/post/", controllers.RenderPostPage)
	http.HandleFunc("/create_post/", func(w http.ResponseWriter, r *http.Request) {
		_, valid := api.ValidateUser(w, r)
		if !valid {
			w.Write([]byte("Unauthorize access"))
			return
		}
		controllers.RenderCreatePostPage(w, r)
	})

	fmt.Println("\033[1;32mServer listening on port http://localhost:8080 ...\033[0m")
	log.Fatal(http.ListenAndServe(":8080", nil))

	if err := funcs.DB.Close(); err != nil {
		fmt.Println("Error :", err)
	}
}

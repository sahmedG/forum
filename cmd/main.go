package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"forum/api"
	"forum/controllers"
	"forum/pkgs/funcs"
)

const (
	limit    = 100             // Max number of requests
	interval = 1 * time.Second // Time interval for rate limiting
)

var (
	tokens = make(chan struct{}, limit)
)

func LoadConfigFromFile(filename string) (*api.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config api.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func init() {
	// Fill the token bucket initially
	go func() {
		ticker := time.NewTicker(interval / limit)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tokens <- struct{}{}:
				default:
				}
			}
		}
	}()
}

func main() {
	funcs.Init()

	// Create a file server to serve static files (CSS, JS, images, etc.)
	fs := http.FileServer(http.Dir("static"))

	// Handle requests for files in the "/static/" path
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	/********************* API endpoints ************************/
	http.HandleFunc("/signup", handleWithRateLimit(api.SignUp)) // Handle signup
	http.HandleFunc("/login", handleWithRateLimit(api.LogIn))   // Handle login
	http.HandleFunc("/logout", handleWithRateLimit(api.LogOut)) // Handle logout

	//authentication part start
	http.HandleFunc("/google", handleWithRateLimit(api.HandleGoogleLogin))             //Handle google
	http.HandleFunc("/google/callback", handleWithRateLimit(api.HandleGoogleCallback)) //Handle google callback
	http.HandleFunc("/github", handleWithRateLimit(api.HandleGithubLogin))             //Handle github
	http.HandleFunc("/github/callback", handleWithRateLimit(api.HandleGithubCallback)) //Handle github callback
	//authentication part end

	http.HandleFunc("/api/islogged", handleWithRateLimit(api.Serve_is_logged))                // Check if user logged in or not
	http.HandleFunc("/api/created_by_user", handleWithRateLimit(api.ByUser_filter_handler))   // posts filtering ex: /api/userposts (when the user is loggedin)
	http.HandleFunc("/api/liked_by_user", handleWithRateLimit(api.ByUser_filter_handler))     // posts filtering ex: /api/userposts (when the user is loggedin)
	http.HandleFunc("/api/create_post", handleWithRateLimit(api.Create_Post))                 // create post
	http.HandleFunc("/api/posts", handleWithRateLimit(api.GetPostsHandler))                   // Retrive posts as JSON
	http.HandleFunc("/api/post/", handleWithRateLimit(api.Get_post_handler))                  // Retrive one post ex: /post/2
	http.HandleFunc("/api/likes_post", handleWithRateLimit(api.LikesPostHandler))             // Handle Likes & Dislikes for Posts
	http.HandleFunc("/api/postlikes", handleWithRateLimit(api.Serve_post_likes_dislikes))     // Serve post likes & dislikes
	http.HandleFunc("/api/add_comment", handleWithRateLimit(api.AddCommentHandler))           // Create comment
	http.HandleFunc("/api/comments", handleWithRateLimit(api.Serve_comments_handler))         // Serve post comments
	http.HandleFunc("/api/likes_comment", handleWithRateLimit(api.LikesCommentHandler))       // Handle Likes & Dislikes for Comments
	http.HandleFunc("/api/commlikes", handleWithRateLimit(api.Serve_comm_likes_dislikes))     // Serve comment likes & dislikes
	http.HandleFunc("/api/create_category", handleWithRateLimit(api.Create_Category_Handler)) // Create category
	http.HandleFunc("/api/categories", handleWithRateLimit(api.Serve_categories_handler))     // Serve categories
	http.HandleFunc("/api/category/", handleWithRateLimit(api.Categories_filter_handler))     // Category filtering ex: /api/category/Technology

	/********************* END ************************/

	// Render pages
	http.HandleFunc("/", handleWithRateLimit(controllers.RenderPage))
	http.HandleFunc("/post/", handleWithRateLimit(controllers.RenderPostPage))
	http.HandleFunc("/create_post/", handleWithRateLimit(func(w http.ResponseWriter, r *http.Request) {
		_, valid := api.ValidateUser(w, r)
		if !valid {
			w.Write([]byte("Unauthorized access"))
			return
		}
		controllers.RenderCreatePostPage(w, r)
	}))

	config, err := LoadConfigFromFile("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	// TLS configuration
	tlsConfig := &tls.Config{}
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		log.Fatal("Failed to load certificate or key file:", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	readTimeout, err := time.ParseDuration(config.ReadTimeout)
	if err != nil {
		log.Fatal("Failed to parse read timeout:", err)
	}
	writeTimeout, err := time.ParseDuration(config.WriteTimeout)
	if err != nil {
		log.Fatal("Failed to parse write timeout:", err)
	}
	idleTimeout, err := time.ParseDuration(config.IdleTimeout)
	if err != nil {
		log.Fatal("Failed to parse idle timeout:", err)
	}
	// Create HTTP server
	server := &http.Server{
		Addr:         config.ServerAddr,
		Handler:      nil,
		TLSConfig:    tlsConfig,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	fmt.Printf("Server listening on https://localhost%s ...\n", config.ServerAddr)
	// log.Printf("Server listening on %s ...\n", config.ServerAddr)
	log.Fatal(server.ListenAndServeTLS(config.CertFile, config.KeyFile))
}

func handleWithRateLimit(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-tokens:
			handler(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	}
}

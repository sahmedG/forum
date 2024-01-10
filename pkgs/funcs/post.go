package funcs

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func GetPostID(postID string) (int, error) {
	var postExists bool
	Id, err := strconv.Atoi(postID)
	if err != nil {
		return 0, err
	}
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM posts WHERE p_id = ?)"
	if err := DB.QueryRow(query, Id).Scan(&postExists); err != nil {
		return 0, err
	}
	if !postExists {
		return 0, fmt.Errorf("post does not exist")
	}

	return Id, nil
}

func CreatePost(userID int, title string, categories []string, content string) error {

	// Check if catigories are existed first
	catIDs := make([]int, 0)
	for _, category := range categories {

		// Fetching Category ID
		catID, err := GetCategoryID(category)
		if err != nil {
			return fmt.Errorf("category does not exist")
		}
		catIDs = append(catIDs, catID)
	}
	// Trimming whitespace from the content
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("message cannot be empty")
	}

	// Inserting post data into the database
	query := "INSERT INTO posts (user_id, title, post) VALUES (?, ?, ?)"
	result, err := DB.Exec(query, userID, title, content)
	if err != nil {
		return fmt.Errorf("failed to insert the post")
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve the last inserted ID")
	}
	postID := int(lastID)

	// Inserting the post category into the database
	for _, catID := range catIDs {
		query = "INSERT INTO threads (post_id, cat_id) VALUES (?, ?)"
		if _, err := DB.Exec(query, postID, catID); err != nil {
			return fmt.Errorf("failed to insert the post category")
		}
	}

	// If no catigories in the post we Inserte it onto General category
	if len(catIDs) == 0 {
		query = "INSERT INTO threads (post_id, cat_id) VALUES (?, ?)"
		if _, err := DB.Exec(query, postID, 1); err != nil {
			return fmt.Errorf("failed to insert the post General category")
		}
	}
	return nil

}

// Func to get posts from database
func Get_posts_from_db() ([]Post_json, error) {
	// Query the database
	rows, err := DB.Query(`
	SELECT user_id, user_profile.user_name, posts.creation_date, posts.title, posts.p_id
	FROM posts
	INNER JOIN user_profile ON posts.user_id = user_profile.user_account_id
	ORDER BY p_id DESC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	results := make([]Post_json, 0)

	// Iterate through the rows
	for rows.Next() {
		var userID, userName, creationDate, title, PostID string
		if err := rows.Scan(&userID, &userName, &creationDate, &title, &PostID); err != nil {
			log.Fatal(err)
		}
		// Retrieve categories for the post
		categories, err := Get_Post_Categories(PostID)
		if err != nil {
			return results, err
		}
		date := creationDate[:10]
		// Do something with the data, for example, add it to the result slice
		post_ideee, _ := strconv.Atoi(PostID)       // converts post_id to integer
		countLikes, _ := CountPostLikes(post_ideee) // gets likes count for this post in this idx
		results = append(results, Post_json{        // Append this post into posts_arr array
			User_ID:       userID,
			User_Name:     userName,
			Creation_Date: date,
			Title:         title,
			PostLikes:     countLikes.LikeCount,
			PostDisLikes:  countLikes.DislikeCount,
			Post_ID:       PostID,
			Category:      categories,
		})

	}

	return results, err
}

func Get_Post(postID string) (Post_json, error) {
	var postDetails Post_json

	// Create the SQL query
	query := `
        SELECT user_id, user_profile.user_name, posts.creation_date, posts.title, posts.post
        FROM posts
        JOIN user_profile ON posts.user_id = user_profile.user_account_id
        WHERE posts.p_id = ?
    `
	// Execute the query and retrieve the row
	row := DB.QueryRow(query, postID)

	// Scan the row values into the postDetails struct
	if err := row.Scan(&postDetails.User_ID, &postDetails.User_Name, &postDetails.Creation_Date, &postDetails.Title, &postDetails.Text); err != nil {
		if err == sql.ErrNoRows {
			// Post not found
			return Post_json{}, fmt.Errorf("post not found")
		}
		return Post_json{}, err
	}

	// Retrieve categories for the post
	categories, err := Get_Post_Categories(postID)
	if err != nil {
		return Post_json{}, err
	}
	postDetails.Creation_Date = postDetails.Creation_Date[0:10]
	postDetails.Category = categories
	post_ideee, _ := strconv.Atoi(postID)
	countLikes, _ := CountPostLikes(post_ideee)
	postDetails.PostLikes = countLikes.LikeCount
	postDetails.PostDisLikes = countLikes.DislikeCount

	// Post details retrieved successfully
	return postDetails, nil
}

func Get_Post_Categories(postID string) ([]string, error) {
	query := `
        SELECT category
        FROM threads
        JOIN category ON threads.cat_id = category.cat_id
        WHERE post_id = ?
    `

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func Get_posts_of_category(postIDs []int) ([]Post_json, error) {
	// Generate the placeholders for the SQL query
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))
	for i := range postIDs {
		placeholders[i] = "?"
		args[i] = postIDs[i]
	}

	// Join placeholders with commas
	inClause := strings.Join(placeholders, ",")

	// Query the database
	query := fmt.Sprintf(`
	SELECT user_id, user_profile.user_name, posts.creation_date, posts.title, posts.p_id
	FROM posts
	INNER JOIN user_profile ON posts.user_id = user_profile.user_account_id
	WHERE posts.p_id IN (%s) 
	ORDER BY p_id DESC`, inClause)

	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to hold the results
	results := make([]Post_json, 0)

	// Iterate through the rows
	for rows.Next() {
		var userID, userName, creationDate, title, postID string
		if err := rows.Scan(&userID, &userName, &creationDate, &title, &postID); err != nil {
			log.Fatal(err)
		}
		// Retrieve categories for the post
		categories, err := Get_Post_Categories(postID)
		if err != nil {
			return results, err
		}
		date := creationDate[:10]

		// Do something with the data, for example, add it to the result slice
		post_ideee, _ := strconv.Atoi(postID)       // converts post_id to integer
		countLikes, _ := CountPostLikes(post_ideee) // gets likes count for this post in this idx
		results = append(results, Post_json{        // Append this post into posts_arr array
			User_ID:       userID,
			User_Name:     userName,
			Creation_Date: date,
			Title:         title,
			PostLikes:     countLikes.LikeCount,
			PostDisLikes:  countLikes.DislikeCount,
			Post_ID:       postID,
			Category:      categories,
		})
	}
	return results, err
}

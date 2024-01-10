package funcs

import (
	"fmt"
	"strings"
)

func CreateComment(userID int, postID int, content string) error {
	// Trimming whitespace from the content
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("message cannot be empty")
	}
	var postExists bool
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM posts WHERE p_id = ?)"
	if err := DB.QueryRow(query, postID).Scan(&postExists); err != nil {
		return fmt.Errorf("post does not exist")
	}
	if !postExists {
		return fmt.Errorf("post does not exist")
	}

	// Inserting the comment data into the database
	query = "INSERT INTO comments (post_id,user_id,comment) VALUES (?,?,?);"
	if _, err := DB.Exec(query, postID, userID, content); err != nil {
		return fmt.Errorf("failed to insert the comment")
	}
	return nil

}

// holds all the Comment info
type CommentResults struct {
	UserID          int    `json:"user_id"`
	UserName        string `json:"user_name"`
	CommentID       int    `json:"comment_id"`
	Comment         string `json:"comment"`
	CreationDate    string `json:"creation_date"`
	Edited          bool   `json:"edited"`           // can be used later to show it the post been edited
	CommentLikes    int    `json:"comment_likes"`    // can be fed from funcs.CountCommentLikes()
	CommentDislikes int    `json:"comment_dislikes"` // can be fed from funcs.CountCommentLikes()
	IsLiked         int    `json:"isLiked"`
}

// Func to get Comment from database
func GetComment(postID int) ([]CommentResults, error) {

	// Query the database
	rows, err := DB.Query(`
	SELECT user_id, user_profile.user_name, comm_id, comment_date, comment, edited
	FROM comments 
	INNER JOIN user_profile ON comments.user_id = user_profile.user_account_id
	WHERE post_id = ?`, postID)

	if err != nil {
		return []CommentResults{}, err
	}
	defer rows.Close()

	// Create a slice to hold the results
	var result []CommentResults

	// Iterate through the rows
	for rows.Next() {
		// Create a struct to hold the current result
		var comment CommentResults

		// Scan the values into the struct fields
		if err := rows.Scan(&comment.UserID, &comment.UserName, &comment.CommentID, &comment.CreationDate, &comment.Comment, &comment.Edited); err != nil {
			return []CommentResults{}, err
		}
		// can be removed and done somewhere else
		LikesCount, _ := CountCommentLikes(comment.CommentID)
		comment.CommentLikes = LikesCount.LikeCount
		comment.CommentDislikes = LikesCount.DislikeCount
		comment.CreationDate = comment.CreationDate[:10]

		// Append the current result to the slice
		result = append(result, comment)
	}

	// Check if any rows were returned
	if len(result) == 0 {
		// No comments found for the given postID
		return []CommentResults{}, fmt.Errorf("no comments found for post ID %d", postID)
	}

	return result, nil
}

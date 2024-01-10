package funcs

import (
	"fmt"
	"log"
)

func PostLikes(userID int, postID int, action int) error {
	var postExists bool
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM posts WHERE p_id = ?)"
	if err := DB.QueryRow(query, postID).Scan(&postExists); err != nil {
		return err
	}
	if !postExists {
		return fmt.Errorf("post does not exist")
	}

	var user_interaction_exists bool
	// Checking if the user interacts in the database
	query = "SELECT EXISTS (SELECT actions_type FROM posts_interaction WHERE post_id = ? AND user_id = ?)"
	if err := DB.QueryRow(query, postID, userID).Scan(&user_interaction_exists); err != nil {
		return err
	}

	// Remove action from database when action is ZERO
	if action == 0 {
		query = "DELETE FROM posts_interaction WHERE post_id = ? AND user_id = ?"
		if _, err := DB.Exec(query, postID, userID); err != nil {
			return fmt.Errorf("failed to delete the Like/Dislike action")
		}
	} else if action == 1 || action == -1 {
		if user_interaction_exists {
			query = "UPDATE posts_interaction SET actions_type = ? WHERE post_id = ? AND user_id = ?"
			if _, err := DB.Exec(query, action, postID, userID); err != nil {
				return fmt.Errorf("failed to insert the Like/Dislike action")
			}
		} else {
			// Inserting the action (like/dislike) data into the database
			query = "INSERT INTO posts_interaction (post_id, user_id, actions_type) VALUES (?, ?, ?)"
			if _, err := DB.Exec(query, postID, userID, action); err != nil {
				return fmt.Errorf("failed to insert the Like/Dislike action")
			}
		}
	} else {
		return fmt.Errorf("actions must be (1, 0, -1) only")
	}

	return nil
}

func CommentLikes(userID int, commentID int, action int) error {
	var commentExists bool
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM comments WHERE comm_id = ?)"
	if err := DB.QueryRow(query, commentID).Scan(&commentExists); err != nil {
		return err
	}
	if !commentExists {
		return fmt.Errorf("comment does not exist")
	}

	var user_interaction_exists bool
	// Checking if the user interacts in the database
	query = "SELECT EXISTS (SELECT actions_type FROM comments_interactions WHERE comment_id = ? AND user_id = ?)"
	if err := DB.QueryRow(query, commentID, userID).Scan(&user_interaction_exists); err != nil {
		return err
	}

	// Remove action from database when action is ZERO
	if action == 0 {
		query = "DELETE FROM comments_interactions WHERE comment_id = ? AND user_id = ?"
		if _, err := DB.Exec(query, commentID, userID); err != nil {
			return fmt.Errorf("failed to delete the Like/Dislike action")
		}
	} else if action == 1 || action == -1 {
		if user_interaction_exists {
			query = "UPDATE comments_interactions SET actions_type = ? WHERE comment_id = ? AND user_id = ?"
			if _, err := DB.Exec(query, action, commentID, userID); err != nil {
				return fmt.Errorf("failed to insert the Like/Dislike action")
			}
		} else {
			// Inserting the action (like/dislike) data into the database
			query = "INSERT INTO comments_interactions (comment_id, user_id, actions_type) VALUES (?, ?, ?)"
			if _, err := DB.Exec(query, commentID, userID, action); err != nil {
				return fmt.Errorf("failed to insert the Like/Dislike action")
			}
		}
	} else {
		return fmt.Errorf("actions must be (1, 0, -1) only")
	}

	return nil

}

type LikeCounts struct {
	LikeCount    int `json:"like_count"`
	DislikeCount int `json:"dislike_count"`
}

func CountPostLikes(postID int) (LikeCounts, error) {
	var postExists bool
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM posts WHERE p_id = ?)"
	if err := DB.QueryRow(query, postID).Scan(&postExists); err != nil {
		return LikeCounts{}, err
	}
	if !postExists {
		return LikeCounts{}, fmt.Errorf("post does not exist")
	}

	var counts LikeCounts

	query = "SELECT actions_type, COUNT(*) as count FROM posts_interaction WHERE post_id = ? GROUP BY actions_type;"
	rows, err := DB.Query(query, postID)
	if err != nil {
		return LikeCounts{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var actionType bool
		var count int
		if err := rows.Scan(&actionType, &count); err != nil {
			return LikeCounts{}, err
		}
		if actionType {
			counts.LikeCount = count
		} else {
			counts.DislikeCount = count
		}
	}

	return counts, nil
}

func CountCommentLikes(commentID int) (LikeCounts, error) {

	var commentExists bool
	// Checking if the post ID exists in the database
	query := "SELECT EXISTS (SELECT 1 FROM comments WHERE comm_id = ?)"
	if err := DB.QueryRow(query, commentID).Scan(&commentExists); err != nil {
		return LikeCounts{}, err
	}
	if !commentExists {
		return LikeCounts{}, fmt.Errorf("comment does not exist")
	}

	var counts LikeCounts

	query = "SELECT actions_type, COUNT(*) as count FROM comments_interactions WHERE comment_id = ? GROUP BY actions_type;"
	rows, err := DB.Query(query, commentID)
	if err != nil {
		return LikeCounts{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var actionType bool
		var count int
		if err := rows.Scan(&actionType, &count); err != nil {
			return LikeCounts{}, err
		}
		if actionType {
			counts.LikeCount = count
		} else {
			counts.DislikeCount = count
		}
	}

	return counts, nil
}

/*
* This function will check if user have an action on specific post, and what type of action it got.
* returns 1 when there is like or -1 when disliked or 0 when there is no action
 */
func Post_is_liked_by_user(act_id int, postID string) (int, error) {
	/********************* Below routine is for checking if there exist like or not ***************************/
	var like_exist bool

	query := "SELECT EXISTS (SELECT actions_type FROM posts_interaction WHERE user_id = ? AND post_id = ?)"
	if err := DB.QueryRow(query, act_id, postID).Scan(&like_exist); err != nil {
		return 0, err
	}
	if !like_exist {
		return 0, fmt.Errorf("no like or dislike")
	}
	/********************* END ***************************/

	/********************* Get the action ***************************/
	// Your SQL query
	query = "SELECT actions_type FROM posts_interaction WHERE user_id = ? AND post_id = ?"

	// Execute the query
	rows, err := DB.Query(query, act_id, postID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var column1 bool // Replace with the actual column types in your table
	// Iterate over the result set
	for rows.Next() {
		if err := rows.Scan(&column1); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Column1: %t\n", column1) for debugging
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	/********************* END ***************************/

	// If action is true the return 1, or else return -1
	if column1 {
		return 1, err
	}
	return -1, err
}

/*
* This function will check if user have an action on specific Comment, and what type of action it got.
* returns 1 when there is like or -1 when disliked or 0 when there is no action
 */
func Comment_is_liked_by_user(act_id int, comm_id string) (int, error) {
	/********************* Below routine is for checking if there exist like or not ***************************/
	var like_exist bool

	query := "SELECT EXISTS (SELECT actions_type FROM comments_interactions WHERE user_id = ? AND comment_id = ?)"
	if err := DB.QueryRow(query, act_id, comm_id).Scan(&like_exist); err != nil {
		return 0, err
	}
	if !like_exist {
		return 0, fmt.Errorf("no like or dislike")
	}
	/********************* END ***************************/

	/********************* Get the action ***************************/
	// Your SQL query
	query = "SELECT actions_type FROM comments_interactions WHERE user_id = ? AND comment_id = ?"

	// Execute the query
	rows, err := DB.Query(query, act_id, comm_id)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var column1 bool // Replace with the actual column types in your table
	// Iterate over the result set
	for rows.Next() {
		if err := rows.Scan(&column1); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Column1: %t\n", column1) for debugging
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	/********************* END ***************************/

	// If action is true the return 1, or else return -1
	if column1 {
		return 1, err
	}
	return -1, err
}

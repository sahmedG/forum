package funcs

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func AddUser(userName string, userEmail string, pwd string) error {
	userTypeID, err := UserTypeID("user")
	if err != nil {
		return err
	}
	//    query := "INSERT INTO users (user_email, user_pwd, user_type) VALUES (?, ?, ?)"

	query := "INSERT INTO users (user_email, user_pwd, user_type) VALUES (?, ?, ?)"
	if _, err := DB.Exec(query, userEmail, pwd, userTypeID); err != nil {

		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("user already exist")
		}
		return err
	}

	err = CreateUserProfile(userName, userEmail)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func CreateUserProfile(userName string, userEmail string) error {
	userID, err := SelectUserID(userEmail)
	if err != nil {
		fmt.Println(err)
	}
	query := "INSERT INTO user_profile (user_account_id, user_name) VALUES (?, ?)"

	if _, err := DB.Exec(query, userID, userName); err != nil {

		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("user already exist")
		}
		return err
	}
	return nil
}

func UserTypeID(userType string) (int, error) {
	userTypeIDQuery := "SELECT uty_id FROM user_type WHERE user_type = ?"
	var userTypeID int
	if err := DB.QueryRow(userTypeIDQuery, userType).Scan(&userTypeID); err != nil {
		return 0, err
	}
	return userTypeID, nil
}

// Retrive userID from database
func SelectUserID(userEmail string) (int, error) {
	query := "SELECT u_id FROM users WHERE user_email = ?"
	row := DB.QueryRow(query, userEmail)

	var userID int
	if err := row.Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("incorrect userEmail / password")
		}
		return 0, err
	}
	return userID, nil
}

// this function  changes the usertype to any other type that is in the database
func ChangeUserType(adminID, userID int, typToBe string) error {
	if !UserIsType(adminID, "admin") {
		return fmt.Errorf("only admins")
	}
	typeID, err := UserTypeID(typToBe)
	if err != nil {
		return err
	}
	if UserIsType(userID, typToBe) {
		return fmt.Errorf("user is already a %s", typToBe)
	}
	query := "UPDATE users SET user_type = ? WHERE u_id = ?"
	if _, err = DB.Exec(query, typeID, userID); err != nil {
		return err
	}
	return nil
}

func UserIsType(userID int, typ string) bool {
	query := "SELECT user_type FROM users WHERE u_id = ?"
	var userType int
	if err := DB.QueryRow(query, userID).Scan(&userType); err != nil {
		return false
	}

	typeID, err := UserTypeID(typ)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return userType == typeID
}

// Gets user hashed password
func GetUserHash(userID int) string {
	query := "SELECT user_pwd FROM users WHERE u_id = ?"
	var userHash string
	if err := DB.QueryRow(query, userID).Scan(&userHash); err != nil {
		return "nil"
	}

	return userHash
}

// Func to get all post IDs of a user from database
func GetUserPosts(userID int, onlyLiked bool) ([]int, error) {
	query := "SELECT p_id FROM posts WHERE user_id = ?"
	if onlyLiked {
		query="SELECT post_id FROM posts_interaction WHERE user_id = ? AND actions_type = 1"
	}
	// Query the database
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the result
	var result []int

	// Iterate through the rows
	for rows.Next() {
		var postID int

		// Scan the values into the variable's address
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}

		// Append the current postID to the result slice
		result = append(result, postID)
	}

	// Check if any rows were returned
	if len(result) == 0 {
		// No posts found for the given userID
		return nil, fmt.Errorf("no posts found for user ID %d", userID)
	}
	return result, nil
}

// func CreateUserType(userID int, userType string) error {
// 	if !UserIsAdmin(userID) {
// 		return fmt.Errorf("user is not allowed to create a type")
// 	}
// 	query := "INSERT INTO user_type (user_type) VALUES (?)"
// 	if _, err := DB.Exec(query, userType); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func UserToMod(adminID,userID int) error{
// 	if !UserIsAdmin(adminID) {
// 		return fmt.Errorf("only admins")
// 	}

// 	if UserIsMod(userID) {
// 		return fmt.Errorf("user is already a moderator")
// 	}

// 	modTypeID, err := UserTypeID("moderator")
// 	if err != nil {
// 		return err
// 	}

// 	query := "UPDATE users SET user_type = ? WHERE u_id = ?"
// 	if _, err = DB.Exec(query, modTypeID, userID); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func UserToAdmin(adminID,userID int) error{
// 	if !UserIsAdmin(adminID) {
// 		return fmt.Errorf("only admins")
// 	}

// 	if UserIsAdmin(userID) {
// 		return fmt.Errorf("user is already an admin")
// 	}

// 	adminTypeID, err := UserTypeID("admin")
// 	if err != nil {
// 		return err
// 	}

// 	query := "UPDATE users SET user_type = ? WHERE u_id = ?"
// 	if _, err = DB.Exec(query, adminTypeID, userID); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func UserIsMod(userID int) bool {
// 	query := "SELECT user_type FROM users WHERE u_id = ?"
// 	var userType int
// 	if err := DB.QueryRow(query, userID).Scan(&userType); err != nil {
// 		return false
// 	}

// 	modTypeID, err := UserTypeID("moderator")
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return userType == modTypeID
// }

// func UserIsAdmin(userID int) bool {
// 	query := "SELECT user_type FROM users WHERE u_id = ?"
// 	var userType int
// 	if err := DB.QueryRow(query, userID).Scan(&userType); err != nil {
// 		return false
// 	}

// 	adminTypeID, err := UserTypeID("admin")
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return userType == adminTypeID
// }

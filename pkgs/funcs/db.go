package funcs

import (
	"database/sql"
	"fmt"
	"forum/pkgs/hashing"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

const DBPath string = "forum.db"

var DB *sql.DB

func Init() {
	// connect to the database
	var New_DB bool
	_, err := os.Stat(DBPath)
	if os.IsNotExist(err) {
		InitateDB(DBPath)
		New_DB = true
	}

	DB, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}

	if New_DB {
		//add admin accunt to the new Dara Base
		CreateAdmin(DB)
	}
}

func InitateDB(DBPath string) {
	DB, err := sql.Open("sqlite3", "forum.db")

	if err != nil {
		log.Fatal(err)
	}
	err = CreateUserTables(DB)
	if err != nil {
		log.Fatal(err)
	}
	err = CreatePostsTables(DB)
	if err != nil {
		log.Fatal(err)
	}
	err = CreateCommentsTables(DB)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateAdmin(DB *sql.DB) error {
	HashedPassword, _ := hashing.HashPassword("admin")
	email := "admin"
	userTypeID, err := UserTypeID("admin")
	if err != nil {
		return err
	}

	query := "INSERT INTO users (user_email, user_pwd, user_type) VALUES (?, ?, ?)"
	if _, err := DB.Exec(query, email, HashedPassword, userTypeID); err != nil {

		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("user already exist")
		}
		return err
	}

	err = CreateUserProfile(email, email)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

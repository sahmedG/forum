package funcs

import (
	"database/sql"
)

func CreateUserTables(db *sql.DB) error {
	//Create users table
	_, err := db.Exec(`	CREATE TABLE users (
			u_id                 INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
			user_email           VARCHAR(50) NOT NULL    ,
			user_pwd             VARCHAR NOT NULL    ,
			user_type            INTEGER NOT NULL    ,
			CURRENT_DATE         TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
			CONSTRAINT unq_users_user_id UNIQUE ( user_email ),
			FOREIGN KEY ( user_type ) REFERENCES user_type( uty_id )
		 )`)
	if err != nil {
		return err
	}

	//Create reports table
	_, err = db.Exec(`CREATE TABLE reports (
		rep_id               INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		mod_id               INTEGER NOT NULL    ,
		post_id              INTEGER NOT NULL    ,
		report               VARCHAR NOT NULL    ,
		report_date          TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
		FOREIGN KEY ( mod_id ) REFERENCES users( u_id )  ,
		FOREIGN KEY ( post_id ) REFERENCES posts( p_id )
	 )`)
	if err != nil {
		return err
	}

	//Create requests table
	_, err = db.Exec(`CREATE TABLE request (
		req_id               INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		user_id              INTEGER NOT NULL    ,
		request_status       CHAR NOT NULL    ,
		creation_date        TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
		UPDATE_status_date   TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
		FOREIGN KEY ( user_id ) REFERENCES users( u_id )
	 )`)
	if err != nil {
		return err
	}

	//Create user_type table
	_, err = db.Exec(`CREATE TABLE user_type (
		uty_id               INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		user_type            CHAR NOT NULL
	 )`)
	if err != nil {
		return err
	}

	//Create user_profile table
	_, err = db.Exec(`CREATE TABLE user_profile (
		profile_id           INTEGER NOT NULL  PRIMARY KEY  ,
		user_account_id      INTEGER NOT NULL    ,
		user_name            VARCHAR(50) NOT NULL    ,
		DOB                  DATE     ,
		first_name           VARCHAR(50)     ,
		last_name            VARCHAR(50)     ,
		creation_date        DATE  DEFAULT CURRENT_DATE   ,
		FOREIGN KEY ( user_account_id ) REFERENCES users( u_id )
	 )`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO user_type (user_type) VALUES ('admin')`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO user_type (user_type) VALUES ('mod')`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO user_type (user_type) VALUES ('user')`)
	if err != nil {
		return err
	}

	return nil
}

func CreatePostsTables(db *sql.DB) error {

	//Cearte Table Posts
	_, err := db.Exec(`CREATE TABLE posts (
		p_id                 INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		user_id              INTEGER NOT NULL    ,
		creation_date        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP   ,
		post                 VARCHAR(150) NOT NULL    ,
		edited               BOOLEAN NOT NULL DEFAULT FALSE   ,
		title                CHAR(150)  DEFAULT 'NULL'   ,
		FOREIGN KEY ( user_id ) REFERENCES users( u_id )
	 );
	`)
	if err != nil {
		return err
	}

	//Cearte Table Categories
	_, err = db.Exec(`CREATE TABLE category (
			cat_id               INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
			category             CHAR NOT NULL
		 );
		`)
	if err != nil {
		return err
	}

	//Cearte Table Threads
	_, err = db.Exec(`CREATE TABLE threads (
		ID                   INTEGER NOT NULL  PRIMARY KEY  ,
		post_id              INTEGER NOT NULL    ,
		cat_id               INTEGER NOT NULL DEFAULT 1   ,
		FOREIGN KEY ( post_id ) REFERENCES posts( p_id )  ,
		FOREIGN KEY ( cat_id ) REFERENCES category( cat_id )
	 );
	`)
	if err != nil {
		return err
	}

	//Cearte Table deleted_posts
	_, err = db.Exec(`CREATE TABLE deleted_posts (
			id                   INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
			post_id              INTEGER NOT NULL    ,
			user_id              INTEGER NOT NULL    ,
			post                 TEXT     ,
			reason               TEXT NOT NULL    ,
			deletion_date        DATE  DEFAULT CURRENT_DATE   ,
			FOREIGN KEY ( post_id ) REFERENCES posts( p_id )
		 );
		`)
	if err != nil {
		return err
	}

	//Cearte Table static_files
	_, err = db.Exec(`CREATE TABLE static_files (
		file_id              INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		post_id              INTEGER NOT NULL    ,
		file_path            VARCHAR NOT NULL    ,
		upload_date          TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
		FOREIGN KEY ( post_id ) REFERENCES posts( p_id )
	 );
	`)
	if err != nil {
		return err
	}

	//Cearte Table post_history
	_, err = db.Exec(`CREATE TABLE post_history (
			id                   INTEGER NOT NULL  PRIMARY KEY  ,
			post_id              INTEGER NOT NULL    ,
			post_details         VARCHAR(255) NOT NULL    ,
			edit_date            DATE  DEFAULT CURRENT_DATE   ,
			FOREIGN KEY ( post_id ) REFERENCES posts( p_id )
		 );
		`)
	if err != nil {
		return err
	}

	//Cearte Table posts_interaction
	_, err = db.Exec(`CREATE TABLE posts_interaction (
		act_id               INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		post_id              INTEGER NOT NULL    ,
		user_id              INTEGER NOT NULL    ,
		actions_type         BOOLEAN NOT NULL    ,
		action_date          TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
		FOREIGN KEY ( post_id ) REFERENCES posts( p_id )  ,
		FOREIGN KEY ( user_id ) REFERENCES users( u_id )
	 );`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO category (category) VALUES ('General')`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO category (category) VALUES ('Engineering')`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO category (category) VALUES ('Travel')`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO category (category) VALUES ('Technology')`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO category (category) VALUES ('Mathematics')`)
	if err != nil {
		return err
	}

	return nil
}

func CreateCommentsTables(db *sql.DB) error {

	// Create comments table
	_, err := db.Exec(`CREATE TABLE comments (
		comm_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		comment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		comment VARCHAR(150) NOT NULL,
		edited BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (post_id) REFERENCES posts(p_id),
		FOREIGN KEY (user_id) REFERENCES users(u_id)
	)`)
	if err != nil {
		return err
	}

	// Create comments_interactions table
	_, err = db.Exec(`CREATE TABLE comments_interactions (
		action_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		actions_type BOOLEAN NOT NULL,
		action_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (comment_id) REFERENCES comments(comm_id),
		FOREIGN KEY (user_id) REFERENCES users(u_id)
	)`)
	if err != nil {
		return err
	}

	// Create comment_history table
	_, err = db.Exec(`CREATE TABLE comment_history (
			id INTEGER NOT NULL PRIMARY KEY,
			comment VARCHAR(255) NOT NULL,
			comment_id INTEGER NOT NULL,
			edit_date DATE NOT NULL DEFAULT CURRENT_DATE,
			FOREIGN KEY (comment_id) REFERENCES comments(comm_id)
		)`)
	if err != nil {
		return err
	}

	//Create deleted comments table
	_, err = db.Exec(`CREATE TABLE deleted_comments (
		id                   INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		comm_id              INTEGER NOT NULL    ,
		user_id              INTEGER NOT NULL    ,
		comment              TEXT     ,
		deletion_date        DATE  DEFAULT CURRENT_DATE   ,
		reason               TEXT NOT NULL    ,
		FOREIGN KEY ( comm_id ) REFERENCES comments( comm_id )
	 )`)
	if err != nil {
		return err
	}

	//Create comments notifications table
	_, err = db.Exec(`CREATE TABLE notifications (
		not_id               INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT ,
		post_id              INTEGER NOT NULL    ,
		not_type             INTEGER NOT NULL    ,
		creation_date        TIMESTAMP  DEFAULT CURRENT_TIMESTAMP   ,
		FOREIGN KEY ( post_id ) REFERENCES posts( p_id )  ,
		FOREIGN KEY ( not_type ) REFERENCES comments( comm_id )  ,
		FOREIGN KEY ( not_type ) REFERENCES likes_dislikes( act_id )
	 )`)
	if err != nil {
		return err
	}

	//Create comments interactions table
	_, err = db.Exec(`	CREATE TABLE likes_dislikes (
		act_id               INTEGER
	 )`)
	if err != nil {
		return err
	}

	return nil
}

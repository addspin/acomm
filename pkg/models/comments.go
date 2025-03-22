package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Comment struct {
	ID     int64  `json:"id"`
	NewsID int64  `json:"news_id"`
	Text   string `json:"text"`
}

var DB *sql.DB

// InitDB initializes the database connection and creates the tables if they don't exist
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	// Create comments table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		news_id INTEGER NOT NULL,
		text TEXT NOT NULL
	);
	`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating database table: %v", err)
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// AddComment adds a new comment to the database
func AddComment(newsID int64, text string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO comments(news_id, text) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newsID, text)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetCommentsByNewsID retrieves all comments for a particular news article
func GetCommentsByNewsID(newsID int64) ([]Comment, error) {
	rows, err := DB.Query("SELECT id, news_id, text FROM comments WHERE news_id = ?", newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.NewsID, &comment.Text); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

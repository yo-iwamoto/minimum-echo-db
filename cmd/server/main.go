package main

import (
	"log"

	_ "github.com/gwenn/gosqlite"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var schema = `
CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	content TEXT
);

CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER,
	content TEXT
);
`

type Post struct {
	ID      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

type Comment struct {
	ID      int    `json:"id" db:"id"`
	PostID  int    `json:"post_id" db:"post_id"`
	Content string `json:"content" db:"content"`
}

type PostResponse struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		db, err := sqlx.Connect("sqlite3", "./database.db")
		if err != nil {
			log.Fatal(err)
		}

		db.MustExec(schema)
		tx := db.MustBegin()
		tx.MustExec(
			"INSERT INTO posts (title, content) VALUES (?, ?)",
			"Hello World",
			"This is my first post",
		)
		tx.Commit()

		posts := []Post{}
		db.Select(&posts, "SELECT * FROM posts")

		response := map[string]interface{}{
			"message": "Hello World",
			"posts":   posts,
		}

		return c.JSON(200, response)
	})

	e.GET("/posts/:id", func(c echo.Context) error {
		db, err := sqlx.Connect("sqlite3", "./database.db")
		if err != nil {
			log.Fatal(err)
		}

		db.MustExec(schema)
		db.MustExec("INSERT INTO comments (post_id, content) VALUES (?, ?)", c.Param("id"), "This is my first comment")

		post := Post{}
		db.Get(&post, "SELECT * FROM posts WHERE id = ?", c.Param("id"))

		comments := []Comment{}
		db.Select(&comments, "SELECT * FROM comments WHERE post_id = ?", c.Param("id"))

		response := PostResponse{
			ID:       post.ID,
			Title:    post.Title,
			Content:  post.Content,
			Comments: comments,
		}

		return c.JSON(200, response)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

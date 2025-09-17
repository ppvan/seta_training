package main

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Tags []string

type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Tags      Tags      `json:"tags" db:"tags"` // Custom type for TEXT[]
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}

	// Use pq.Array for PostgreSQL array handling
	var tags pq.StringArray
	if err := tags.Scan(value); err != nil {
		return err
	}

	*t = Tags(tags)
	return nil
}

// Value implements the driver.Valuer interface for writing to database
func (t Tags) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return pq.Array([]string(t)).Value()
}

// String returns a string representation of tags
func (t Tags) String() string {
	return fmt.Sprintf("[%s]", strings.Join([]string(t), ", "))
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the HTTP
	// request body (not that the field names and types in the struct are a subset of the Movie
	// struct). This struct will be our *target decode destination*.
	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	// Use the readJSON() helper to decode the request body into the struct.
	// If this returns an error we send the client the error message along with
	// a 400 Bad Request status code.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := Post{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}

	query := `
	INSERT INTO posts(title, content, tags)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err = app.db.QueryRow(query, post.Title, post.Content, post.Tags).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

type ActivityLog struct {
	ID       int       `json:"id" db:"id"`
	Action   string    `json:"action" db:"action"`
	PostID   int       `json:"post_id" db:"post_id"`
	LoggedAt time.Time `json:"logged_at" db:"logged_at"`
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

func (me *application) InsertPost(p *Post) (*Post, error) {
	tx, err := me.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert post
	postQuery := `
        INSERT INTO posts(title, content, tags)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	err = tx.QueryRow(postQuery, p.Title, p.Content, p.Tags).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert post: %w", err)
	}

	// Insert activity log
	logQuery := `
        INSERT INTO activity_logs(action, post_id, logged_at)
        VALUES ($1, $2, NOW())
    `
	_, err = tx.Exec(logQuery, "new_post", p.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert activity log: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return p, nil
}

func (me *application) FindPostsByTag(tag string) ([]Post, error) {
	query := `SELECT id, title, content, tags, created_at FROM posts WHERE $1 = ANY(tags)`
	var posts []Post
	rows, err := me.db.Query(query, tag)
	if err != nil {
		return nil, fmt.Errorf("error finding posts by tag: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Tags, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		posts = append(posts, post)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return posts, nil
}

var ErrNotFound = fmt.Errorf("resource not found")

func (me *application) GetAndCachePost(id int) (*Post, error) {
	// Cache check using cache aside, TTL 5 mins
	cacheKey := fmt.Sprintf("post:%d", id)
	ctx := context.Background()

	cached, err := me.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var post Post
		if err := json.Unmarshal([]byte(cached), &post); err != nil {
			fmt.Printf("Cache unmarshal error: %v\n", err)
		} else {
			return &post, nil
		}
	}

	// Cache miss or error - get from database
	query := `SELECT id, title, content, tags, created_at FROM posts WHERE id = $1`
	var post Post
	err = me.db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.Tags, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error retrieving post: %w", err)
	}

	// Store in cache for next time (5 minute TTL)
	postJSON, err := json.Marshal(post)
	if err != nil {
		// Log marshal error but still return the post
		fmt.Printf("Cache marshal error: %v\n", err)
	} else {
		err = me.rdb.Set(ctx, cacheKey, postJSON, 5*time.Minute).Err()
		if err != nil {
			// Log cache set error but still return the post
			fmt.Printf("Cache set error: %v\n", err)
		}
	}

	return &post, nil
}

func (me *application) UpdatePost(id int, p Post) (*Post, error) {
	postQuery := `
        UPDATE posts
        SET title = $1, content = $2, tags = $3
        RETURNING id, title, content, tags, created_at
    `
	err := me.db.QueryRow(postQuery, p.Title, p.Content, p.Tags).Scan(&p.ID, &p.Title, &p.Content, &p.Tags, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Clear cache
	cacheKey := fmt.Sprintf("post:%d", id)
	me.rdb.Del(context.Background(), cacheKey)

	return &p, nil
}

func (me *application) SearchPostFullText(query string) ([]Post, error) {
	sqlQuery := `
		SELECT id, title, content, tags, created_at
		FROM posts
		WHERE search_vector @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(search_vector, plainto_tsquery('english', $1)) DESC
	`

	rows, err := me.db.Query(sqlQuery, query)
	if err != nil {
		return nil, fmt.Errorf("error executing full-text search: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Tags, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning search result: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return posts, nil
}

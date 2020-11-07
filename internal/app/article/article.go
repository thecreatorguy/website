// Package article implements the article type for manipulating from the database as well as REST API routes
// for communicating over HTTP
package article

import (
	"html/template"
	"time"
	"website/internal/sqlconn"
)

// Article represents a single blog post
type Article struct {
	URLKey string `json:"url_key"`
	Title string `json:"title"`
	Content template.HTML `json:"content"`
	ReleaseAt sqlconn.NullTime `json:"release_at"`
	CreatedAt sqlconn.NullTime `json:"created_at"`
	UpdatedAt sqlconn.NullTime `json:"updated_at"`
}

// GetArticle retrieves an Article by its primary key, the URLKey
func GetArticle(key string) (Article, error) {
	row := sqlconn.Pool.QueryRow("SELECT * FROM articles WHERE url_key = $1", key)

	a := Article{}
	err := row.Scan(&a.URLKey, &a.Title, &a.Content, &a.ReleaseAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return a, err
	}

	return a, nil
}

// GetArticles retrieves all articles that exist, sorted by most recently released
func GetArticles() ([]Article, error) {
	rows, err := sqlconn.Pool.Query("SELECT * FROM articles ORDER BY release_at DESC")
	if err != nil {
		return nil, err
	}

	as := []Article{}
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.URLKey, &a.Title, &a.Content, &a.ReleaseAt, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}

		as = append(as, a)
	}
	
	return as, nil
}

// IsReleased finds if an article is released or not
func (a Article) IsReleased() bool {
	return a.ReleaseAt.Valid && a.ReleaseAt.Time.Before(time.Now())
}

// Save saves the representation of the article in the database
func (a Article) Save() error {
	query := `
	INSERT INTO articles (url_key, title, content, release_at, created_at, updated_at) 
	VALUES($1, $2, $3, $4, $5, $6) 
	ON CONFLICT (url_key) DO UPDATE 
		SET title = $2, content = $3, release_at = $4, created_at = $5, updated_at = $6`
	_, err := sqlconn.Pool.Exec(query, a.URLKey, a.Title, a.Content, a.ReleaseAt, a.CreatedAt, a.UpdatedAt)
	return err
}

// Delete deletes the article from the database
func (a Article) Delete() error {
	_, err := sqlconn.Pool.Exec("DELETE FROM articles WHERE url_key = $1", a.URLKey)
	return err
}

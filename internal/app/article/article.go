package article

import (
	"database/sql"
	"website/internal/sqlconn"
)

type Article struct {
	URLKey string `json:"url_key"`
	Title string `json:"title"`
	Content string `json:"content"`
	ReleaseAt sql.NullTime `json:"release_at"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

func GetArticle(key string) (Article, error) {
	row := sqlconn.Pool.QueryRow("SELECT * FROM articles WHERE url_key = $1", key)

	a := Article{}
	err := row.Scan(&a.URLKey, &a.Title, &a.Content, &a.ReleaseAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return a, err
	}

	return a, nil
}

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
package sqlconn

import "github.com/jackc/pgx"

var Pool pgx.ConnPool

func init() {
	// os.Getenv("POSTGRES_URI")
}
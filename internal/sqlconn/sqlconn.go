package sqlconn

import (
	"context"
	"os"

	"github.com/jackc/pgx"
)

var Pool *pgx.ConnPool

func init() {
	conf, err := pgx.ParseConnectionString(os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}
	Pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: conf})
	if err != nil {
		panic(err)
	}

	conn, err := Pool.Acquire()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = conn.Ping(context.Background())
	if err != nil {
		panic(err)
	}
}
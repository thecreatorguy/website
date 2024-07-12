package sqlconn

import (
	"database/sql"
	"os"
	"strings"

	
	_ "github.com/ncruces/go-sqlite3/embed"
	_ "github.com/ncruces/go-sqlite3/driver"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", os.Getenv("SQLITE_DB_FILE"))
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	initFile, err := os.ReadFile(os.Getenv("SQLITE_INIT_FILE"))
	if err != nil {
		panic(err)
	}

	statements := strings.Split(string(initFile), ";\n")

	for _, statement := range statements {
		_, err = DB.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}
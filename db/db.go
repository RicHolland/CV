package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type JobHistory struct {
	Role, Company, FromDate, ToDate, Description string
}

const job_table = `\
CREATE TABLE job_history (
    role text not null,
    company text not null,
    from_date text not null,
    to_date text,
    description text
);`

func connect() *sql.DB {
	db, err := sql.Open("sqlite3", "./cv.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

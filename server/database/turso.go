package database
/*

import (
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var AToken = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDk1NTE4NjAsImlkIjoiOGU5MTRkMDEtZGExYS0xMWVlLWI1ZDQtZjY5YmQyNmE1NGJhIn0.yH-QfskJZoc3Jzq7me6y6MKtPlTcv_3IP7zfFeJbdirUO4O5xyx-0MI1tjwXQ9jAdwTA_0q3nHJA5NJoYyevBw"

type Users struct {
	ID int
	username string
	password string
}

func InitDB() {
	sdf

	url := "libsql://site-sb-fox1n69.turso.io?authToken=AToken"

	db, err := sql.Open("libsql", url)
	if err != nil {
		logrus.Fatal(os.Stderr, "failed connect to database", url)
		os.Exit(1)
	}
	logrus.Infoln("Databases connection...")

	defer db.Close()
}
*/
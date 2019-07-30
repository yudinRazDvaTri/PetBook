package driver

import (
	//"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sqlx.DB {
	var err error
	var db *sqlx.DB
	host := os.Getenv("HOST_POSTGRES")
	port := os.Getenv("PORT_POSTGRES")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sqlx.Open("postgres", connStr)
	logFatal(err)
	err = db.Ping()
	if err != nil {
		logFatal(fmt.Errorf("can't ping, err: %s", err.Error()))
	}
	return db
}

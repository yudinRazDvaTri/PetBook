package driver

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/subosito/gotenv"
	"os"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to open .env file.\n")
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
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to open connection.\n")
	}

	err = db.Ping()
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to ping server.\n")
	}
	fmt.Println("Server started.")
	return db
}

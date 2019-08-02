package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)
func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found:%v", err)
	}
}
func main() {
	var err error
	migrations := getMigrations()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT_BD"),
		os.Getenv("SERVICE_USER"), os.Getenv("SERVICE_PASSWORD"), os.Getenv("SERVICE_DBNAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("can't open db:%v", err)
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Printf("can't work migrations:%v", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
	err = db.Ping()
	if err != nil {
		log.Printf("can't ping:%v", err)
	}

}
func getMigrations() *migrate.FileMigrationSource {
	return &migrate.FileMigrationSource{
		Dir: "migrations",
	}
}

package gomigrations

import (
	//"database/sql"
	//"os"
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found:%v", err)
	}
}

func Migrate(db *sqlx.DB) (err error) {
	migrations := getMigrations()
	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to exec migrations.\n")
	}
	fmt.Printf("Applied %d gomigrations!\n", n)
	err = db.Ping()
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to ping server.\n")
	}
	return
}
func getMigrations() *migrate.FileMigrationSource {
	return &migrate.FileMigrationSource{
		Dir: "gomigrations/migrations",
	}
}

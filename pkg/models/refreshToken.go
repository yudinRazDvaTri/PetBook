package models

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/jmoiron/sqlx"
	"time"
)

type RefreshToken struct {
	Id          int    `db:"id"`
	UserId      int    `db:"user_id"`
	TokenString string `db:"token_string"`
}

type RefreshTokenStore struct {
	DB *sqlx.DB
}

type RefreshTokenStorer interface {
	RefreshTokenExists(userId int, token string) error
	UpdateRefreshToken(userId int, token string, lastUpdateAt time.Time) error
}

func (c *RefreshTokenStore) UpdateRefreshToken(userId int, token string, lastUpdateAt time.Time) error {
	_, err := c.DB.Exec(`INSERT INTO refresh_tokens (user_id, token_string) values ($1, $2)
								ON CONFLICT (user_id) DO UPDATE 
								SET token_string = $2,
								last_update_at = $3;`, userId, token, lastUpdateAt)
	if err != nil {
		return fmt.Errorf("Error occurred while trying to insert/update refresh token in db: %v.\n", err)
	}

	return nil
}

func (c *RefreshTokenStore) RefreshTokenExists(userId int, token string) error {
	var exists bool
	err := c.DB.QueryRow( `SELECT EXISTS
								(SELECT 1 
								FROM refresh_tokens
								WHERE token_string = $1 
								AND user_id = $2)`, token, userId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("Error occurred while trying to select refresh token in db: %v.\n", err)
	}

	if !exists {
		return &utilerr.TokenDoesNotExist{Description: "Refresh token does not exist!"}
	}

	return nil
}
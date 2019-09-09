package models

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	RefreshTokenExists(userId int, token string, userAgent string) error
	UpdateRefreshToken(userId int, token string, lastUpdateAt time.Time, userAgent string) error
	DeleteRefreshToken(token string) error
}

// TODO: implement user-agent and token unique error
func (c *RefreshTokenStore) UpdateRefreshToken(userId int, token string, lastUpdateAt time.Time, userAgent string) error {
	_, err := c.DB.Exec(`INSERT INTO refresh_tokens (user_id, token_string, last_update_at, user_agent) values ($1, $2, $3, $4)
								ON CONFLICT (user_id, user_agent) DO UPDATE
								SET token_string = $2,
								last_update_at = $3`, userId, token, lastUpdateAt, userAgent)
	if err != nil {
		if _, ok := err.(*pq.Error); ok {
			return &utilerr.UniqueTokenError{Description: "Token/user_agent is already in database.\n"}
		}
	}

	if err != nil {
		return fmt.Errorf("Error occurred while trying to insert/update refresh token in db: %v.\n", err)
	}

	return nil
}

func (c *RefreshTokenStore) DeleteRefreshToken(token string) error {
	res, err := c.DB.Exec(`DELETE FROM refresh_tokens
								WHERE token_string = $1`, token)
	if err != nil {
		return fmt.Errorf("Error occurred while trying to delete refresh token in db: %v.\n", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error occurred while trying to count affected rows in users in db: %v.\n", err)
	}

	if num == 0 {
		return &utilerr.TokenDoesNotExist{Description: "Refresh token does not exist!"}
	}

	return nil
}

func (c *RefreshTokenStore) RefreshTokenExists(userId int, token string, userAgent string) error {
	var exists bool
	err := c.DB.QueryRow(`SELECT EXISTS
								(SELECT 1 
								FROM refresh_tokens
								WHERE token_string = $1 
								AND user_id = $2
								AND user_agent = $3)`, token, userId, userAgent).Scan(&exists)
	if err != nil {
		return fmt.Errorf("Error occurred while trying to delete refresh token in db: %v.\n", err)
	}

	if !exists {
		return &utilerr.TokenDoesNotExist{Description: "Refresh token does not exist!"}
	}

	return nil
}

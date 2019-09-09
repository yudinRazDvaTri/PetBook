package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FollowersStore struct {
	DB *sqlx.DB
}
type FollowersStorer interface {
	GetFollowers(userID int) ([]*FollowerPets, error)
	GetFollowing(userID int) (pets []*FollowerPets, err error)
	Followed(userId int, followUser int) error
	UnFollowed(userId int, followUser int) error
}

type FollowerPets struct {
	Name        string `json:"name" db:"name"'`
	Description string `json:"description" db:"description"'`
	UserID      int    `json:"user_id" db:"user_id"`
}

//Function that is called in the template. It returns a boolean value,
//hether the user can subscribe to this user
func (f *FollowerPets) CanFollow(userID int, petsFollowing []*FollowerPets) bool {
	if f.UserID == userID {
		return false
	}
	for _, val := range petsFollowing {
		if val.UserID == f.UserID {
			return false
		}
	}
	return true
}

//Get all subscribers of this user
func (f *FollowersStore) GetFollowers(userID int) (pets []*FollowerPets, err error) {
	rows, err := f.DB.Query("select name, description, p.user_id from pets p inner join followers f on p.user_id=f.user_id where f.follower_id=$1;", userID)
	if err != nil {
		err = fmt.Errorf("Can't read followers from db: %v", err)
		return
	}
	defer rows.Close()
	err = sqlx.StructScan(rows, &pets)
	if err != nil {
		err = fmt.Errorf("Can't scan topics-rows from db: %v", err)
	}
	return
}

//Get users that this user is subscribed to
func (f *FollowersStore) GetFollowing(userID int) (pets []*FollowerPets, err error) {
	rows, err := f.DB.Query("select p.name, p.description, p.user_id from pets p inner join followers f on p.user_id=f.follower_id where f.user_id=$1;", userID)
	if err != nil {
		err = fmt.Errorf("Can't read following from db: %v", err)
		return
	}
	defer rows.Close()
	err = sqlx.StructScan(rows, &pets)
	if err != nil {
		err = fmt.Errorf("Can't scan topics-rows from db: %v", err)
	}
	return
}

//Subscribe to user
func (f *FollowersStore) Followed(userId int, followUser int) error {
	_, err := f.DB.Exec(`insert into followers values ($1,$2);`, userId, followUser)
	if err != nil {
		err = fmt.Errorf("cannot added follower in followers table of db: %v", err)
	}
	return err
}

//Unsubscribe to user
func (f *FollowersStore) UnFollowed(userId int, followUser int) error {
	_, err := f.DB.Exec(`delete from followers where user_id=$1 and follower_id=$2;`, userId, followUser)
	if err != nil {
		err = fmt.Errorf("cannot delete follower in followers table of db: %v", err)
	}
	return err
}

package controllers

import (
	//"fmt"
	"github.com/jmoiron/sqlx"
	"petbook/models"
	"petbook/store"
)

type UserController struct {
	DB *sqlx.DB
}

func (c *UserController) Register(u *models.User) error {
	err := store.Register(c.DB, u)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserController) GetUser(u *models.User, email string) error {
	u.Login = email
	err := store.GetUser(c.DB, u)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserController) ChangePassword(u *models.User, newPassword string) error {
	err := store.ChangePassword(c.DB, u, newPassword)
	if err != nil {
		return err
	}
	u.Password = newPassword
	return nil
}

func (c *UserController) Login(userChecking *models.User) error {
	err := store.Login(c.DB, userChecking)
	if err != nil {
		return err
	}
	return nil
}

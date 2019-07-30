package controllers

import (
	"github.com/jmoiron/sqlx"
	"petbook/models"
	"petbook/store"
)

type UserController struct {
	DB *sqlx.DB
}

func (c *UserController) CreateUser(u *models.User) error {
	err := store.CreateUser(c.DB, u)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserController) GetUser(u *models.User, login string) error {
	u.Login = login
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

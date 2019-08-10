package utilerr

import "fmt"

type UniqueTaken struct {
	Description string
}

type WrongEmail struct {
	Email string
}

type WrongPassword struct {
	Description string
}

func (e *UniqueTaken) Error() string {
	return e.Description
}

func (e *WrongEmail) Error() string {
	return fmt.Sprintf("Account with email '%s' not found.", e.Email)
}

func (e *WrongPassword) Error() string {
	return e.Description
}
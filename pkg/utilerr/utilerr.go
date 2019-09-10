package utilerr

type UniqueTaken struct {
	Description string
}

type WrongCredentials struct {
	Description string
}

type TokenDoesNotExist struct {
	Description string
}

type UniqueTokenError struct {
	Description string
}

type PetDoesNotExist struct {
	Description string
}

func (e *UniqueTaken) Error() string {
	return e.Description
}

func (e *WrongCredentials) Error() string {
	return e.Description
}

func (e *TokenDoesNotExist) Error() string {
	return e.Description
}

func (e *UniqueTokenError) Error() string {
	return e.Description
}

func (e *PetDoesNotExist) Error() string {
	return e.Description
}
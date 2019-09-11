package utilerr

// Provides some basic custom errors, which make error handling easier.

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

type VetDoesNotExist struct {
	Description string
}

type LogoDoesNotExist struct {
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

func (e *LogoDoesNotExist) Error() string {
	return e.Description
}

func (e *VetDoesNotExist) Error() string {
	return e.Description
}

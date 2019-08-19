package controllers

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/models/search"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/dpgolang/PetBook/pkg/view"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
	//"github.com/gorilla/mux"
	//"github.com/gorilla/sessions"
	//"github.com/jmoiron/sqlx"
	//"github.com/subosito/gotenv"
)

type Controller struct {
	UserStore         models.UserStorer
	PetStore          models.PetStorer
	RefreshTokenStore models.RefreshTokenStorer
	ForumStore        forum.ForumStorer
	SearchStore       search.SearchStorer
	BlogStore         models.BlogStorer
}

const (
	patternEmail    = `^\w+@\w+\.\w+$`
	patternPassword = `^.{6,}$`
	patternAnyChar  = `.*\S.*`
	patternOnlyNum  = `^[0-9]*$`
)

func (c *Controller) LoginGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "login")
	}
}
func (c *Controller) LoginPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")
		var userID int
		var err error
		if userID, err = c.UserStore.Login(email, password); err != nil {
			switch e := err.(type) {
			case *utilerr.WrongCredentials:
				// TODO: display flash-message
				fmt.Fprint(w, e.Error())
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			default:
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var tokens authentication.Tokens
		tokens, err = authentication.GenerateTokenPair(userID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "accessToken",
			Value:   tokens.AccessTokenValue,
			Expires: tokens.AccessExpirationTime,
		})

		if err := c.RefreshTokenStore.UpdateRefreshToken(userID, tokens.RefreshTokenValue, tokens.RefreshExpirationTime); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "refreshToken",
			Value:   tokens.RefreshTokenValue,
			Expires: tokens.RefreshExpirationTime,
		})

		_, err = c.UserStore.GetPet(userID)
		if err != nil {
			http.Redirect(w, r, "/petcabinet", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/mypage", http.StatusFound)
	}
}

func (c *Controller) RegisterGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "register")
	}
}

// TODO: reduce repeating code
func (c *Controller) RegisterPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		login := r.FormValue("login")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		//userType := r.FormValue("userType")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")

		if matched, err := regexp.Match(patternAnyChar, []byte(login)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match login.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		if matched, err := regexp.Match(patternEmail, []byte(email)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match email.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		if matched, err := regexp.Match(patternAnyChar, []byte(firstName)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match first name.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		if matched, err := regexp.Match(patternAnyChar, []byte(lastName)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match last name.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		if matched, err := regexp.Match(patternPassword, []byte(password)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match password.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		if password != confirmPassword {
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			logger.Error(err, "Error occurred while trying to hash password.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Email:     email,
			Login:     login,
			Firstname: firstName,
			Lastname:  lastName,
			Password:  string(hashedPassword),
		}

		if err := c.UserStore.Register(&user); err != nil {
			switch e := err.(type) {
			case *utilerr.UniqueTaken:
				// TODO: display flash-message
				fmt.Fprint(w, e.Error())
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			default:
				logger.Error(e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func (c *Controller) LogoutGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		refreshToken, err := r.Cookie("refreshToken")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		refreshTokenString := refreshToken.Value

		if err = c.RefreshTokenStore.DeleteRefreshToken(refreshTokenString); err != nil {
			switch e := err.(type) {
			case *utilerr.TokenDoesNotExist:
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			default:
				logger.Error(e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "accessToken",
			Expires: time.Unix(0, 0),
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "refreshToken",
			Expires: time.Unix(0, 0),
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

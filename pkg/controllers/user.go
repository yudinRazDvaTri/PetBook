package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/models/forum"
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
	UserStore  models.UserStorer
	PetStore   models.PetStorer
	ForumStore forum.ForumStorer
}

const (
	patternEmail    = `^\w+@\w+\.\w+$`
	patternPassword = `^.{6,}$`
	patternAnyChar  = `.*\S.*`
	patternOnlyNum = `^[0-9]*$`
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

		user := models.User{
			Email:    email,
			Password: []byte(password),
		}

		if err := c.UserStore.Login(&user); err != nil {
			switch e := err.(type) {
			case *utilerr.WrongEmail:
				// TODO: display flash-message
				fmt.Fprint(w, e.Error())
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			case *utilerr.WrongPassword:
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

		expirationTime := time.Now().Add(30 * time.Minute)

		claims := &authentication.Claims{
			Id: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

		tokenString, err := token.SignedString(authentication.Keys.SignKey)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err, "Error occurred while trying to sign token.\n")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		_, err = c.UserStore.GetPet(&user)
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
			Password:  hashedPassword,
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

package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/dpgolang/PetBook/pkg/authentication"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/models/forum"
	"github.com/dpgolang/PetBook/pkg/models/search"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/dpgolang/PetBook/pkg/view"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Controller struct {
	UserStore         models.UserStorer
	PetStore          models.PetStorer
	RefreshTokenStore models.RefreshTokenStorer
	ForumStore        forum.ForumStorer
	SearchStore       search.SearchStorer
	BlogStore         models.BlogStorer
	ChatStore         models.ChatStorer
	FollowersStore 	  models.FollowersStorer
	MediaStore        models.MediaStorer
	VetStore  		  models.VetStorer
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
			Path:    "/",
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
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "oauth",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		//_, err = c.UserStore.GetPet(userID)
		//if err != nil {
		//	http.Redirect(w, r, "/petcabinet", http.StatusFound)
		//	return
		//}
		c.cabinetFilled(userID,w,r)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (c *Controller) cabinetFilled(id int,w http.ResponseWriter, r *http.Request) {
	role:=c.UserStore.GetUserRole(id)
	if role=="pet" {
		_, err := c.UserStore.GetPet(id)
		if err != nil {
			http.Redirect(w, r, "/petcabinet", http.StatusFound)
			return
		}
	}else if role =="vet"{
		_, err := c.UserStore.GetVet(id)
		if err != nil {
			http.Redirect(w, r, "/vetcabinet", http.StatusFound)
			return
		}
	}
}




func (c *Controller) LoginGoogleGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState := generateStateOauthCookie(w)
		u := authentication.GoogleOauthConfig.AuthCodeURL(oauthState)+"&access_type=offline"
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(1 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (c *Controller) GoogleCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState, _ := r.Cookie("oauthstate")

		http.SetCookie(w, &http.Cookie{
			Name:    "accessToken",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "refreshToken",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "oauthstate",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		if r.FormValue("state") != oauthState.Value {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		googleToken, err := getGoogleOauthToken(r.FormValue("code"))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		response, err := http.Get(authentication.OauthGoogleUrlAPI + googleToken.AccessToken)
		if err != nil {
			logger.Error("Error occurred while trying to get user info: %v.\n", err.Error())
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("Error occurred while trying to read user info bytes: %v.\n",err)
			return
		}

		var googleUserInfo authentication.GoogleUserInfo

		if err := json.Unmarshal(contents, &googleUserInfo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("Error occurred while trying to unmarshal user info: %v.\n",err)
			return
		}

		userId, err := c.UserStore.LoginOauth(googleUserInfo.Email)

		if err != nil {
			switch e := err.(type) {
			case *utilerr.UniqueTaken:
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

		gob.Register(googleToken)
		value := map[string]interface{} {
			"accessToken": googleToken,
			"userId": userId,
		}

		if encoded, err := authentication.SCookie.Encode("oauth", value); err == nil {
			cookie := &http.Cookie{
				Name:  "oauth",
				Value: encoded,
				// Expiration time of cookie which stores oauth information was set twice as much as google oauth token expiration time.
				// (Google access token expiration time is 3600 seconds)
				Expires: time.Now().Add(7200 * time.Second),
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		} else {
			logger.Error(err.Error(), "; Error occurred while trying to encode cookie.\n", )
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func getGoogleOauthToken(code string) (*oauth2.Token, error) {
	token, err := authentication.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return token, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	//authentication.GoogleOauthConfig.Client(context.Background(), token)
	return token, nil
}

func (c *Controller) RegisterGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles:=c.UserStore.GetUserEnums()
		view.GenerateHTML(w, roles, "register")
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
		userType := r.FormValue("user-role")
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
			Role:userType,
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
		if err == nil {
			refreshTokenString := refreshToken.Value
			if err = c.RefreshTokenStore.DeleteRefreshToken(refreshTokenString); err != nil {
				switch e := err.(type) {
				case *utilerr.TokenDoesNotExist:
				default:
					logger.Error(e)
					http.Error(w, e.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "accessToken",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "refreshToken",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "oauth",
			Expires: time.Unix(0, 0),
			Path:    "/",
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
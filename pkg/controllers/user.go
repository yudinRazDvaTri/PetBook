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
	gorillaContext "github.com/gorilla/context"
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
	FollowersStore    models.FollowersStorer
	MediaStore        models.MediaStorer
	VetStore          models.VetStorer
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

// Controller which handles user logging in.
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

		userAgent := authentication.GetUserAgent(r)

		if err := c.RefreshTokenStore.UpdateRefreshToken(userID, tokens.RefreshTokenValue, tokens.RefreshExpirationTime, userAgent); err != nil {
			switch e := err.(type) {
			case *utilerr.UniqueTokenError:
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			default:
				logger.Error(e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
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

		http.Redirect(w, r, "/mypage", http.StatusSeeOther)
		return
	}
}

func (c *Controller) LoginGoogleGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState := generateStateOauthCookie(w)
		u := authentication.GoogleOauthConfig.AuthCodeURL(oauthState) + "&access_type=offline"
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

// Generating oauth state. State is a token to protect the user from CSRF attacks.
// You must always provide a non-empty string and validate that it matches the
// the state query parameter on your redirect callback.
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(1 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

// This controller is called when the resource owner follows the RedirectURL,
// which was set in GoogleOauthConfig.
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
			logger.Error(err, "Error occurred while trying to get user info.\n")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err, "Error occurred while trying to read user info bytes.\n")
			return
		}

		var googleUserInfo authentication.GoogleUserInfo

		if err := json.Unmarshal(contents, &googleUserInfo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			logger.Error(err, "Error occurred while trying to unmarshal user info.\n")
			return
		}

		userId, err := c.UserStore.LoginOauth(googleUserInfo.Email)

		if err != nil {
			switch e := err.(type) {
			case *utilerr.UniqueTaken:
				// TODO: display flash-message
				fmt.Fprint(w, e.Error())
				return
			default:
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// To store custom types in SecureCookie, they must be registered first using gob.Register().
		gob.Register(googleToken)
		value := map[string]interface{}{
			"accessToken": googleToken,
			"userId":      userId,
		}

		if encoded, err := authentication.SCookie.Encode("oauth", value); err == nil {
			cookie := &http.Cookie{
				Name:  "oauth",
				Value: encoded,
				// Expiration time of cookie which stores oauth information was set twice as much as google oauth token expiration time.
				// (Google access token expiration time is 3600 seconds)
				Expires: time.Now().Add(7200 * time.Second),
				Path:    "/",
			}
			http.SetCookie(w, cookie)
		} else {
			logger.Error(err.Error(), "; Error occurred while trying to encode cookie.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/mypage", http.StatusSeeOther)
		return
	}
}

func getGoogleOauthToken(code string) (*oauth2.Token, error) {
	token, err := authentication.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return token, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	return token, nil
}

func (c *Controller) RegisterGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.GenerateHTML(w, nil, "register")
	}
}

// Controller which handles user registration.
func (c *Controller) RegisterPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		login := r.FormValue("login")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")

		if matched, err := regexp.Match(patternAnyChar, []byte(login)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match login.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if matched, err := regexp.Match(patternEmail, []byte(email)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match email.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if matched, err := regexp.Match(patternPassword, []byte(password)); !matched || err != nil {
			if err != nil {
				logger.Error(err, "Error occurred while trying to match password.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		if password != confirmPassword {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
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
				return
			default:
				logger.Error(e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (c *Controller) RoleGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := c.UserStore.GetUserRole(gorillaContext.Get(r, "id").(int))
		if err != nil {
			logger.Error(err, "Error occurred while trying to user roles enum.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if role != "" {
			http.Redirect(w, r, "/mypage", http.StatusSeeOther)
			return
		}

		roles, err := c.UserStore.GetUserEnums()
		if err != nil {
			logger.Error(err, "Error occurred while trying to user roles enum.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		view.GenerateHTML(w, roles, "role")
	}
}

func (c *Controller) RolePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		userId := gorillaContext.Get(r, "id").(int)

		role, err := c.UserStore.GetUserRole(gorillaContext.Get(r, "id").(int))
		if err != nil {
			logger.Error(err, "Error occurred while trying to user roles enum.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if role != "" {
			http.Redirect(w, r, "/mypage", http.StatusSeeOther)
			return
		}

		roleFormValue := r.FormValue("user-role")
		if roleFormValue != "pet" && roleFormValue != "vet" {
			http.Redirect(w, r, "/role", http.StatusSeeOther)
			return
		}

		if err := c.UserStore.SetUserRole(roleFormValue, userId); err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if roleFormValue == "pet" {
			http.Redirect(w, r, "/petcabinet", http.StatusSeeOther)
			return
		} else if roleFormValue == "vet" {
			http.Redirect(w, r, "/vetcabinet", http.StatusSeeOther)
			return
		}

		http.Error(w, "Wrong roleFormValue!", http.StatusNotFound)
		return

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
					break
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

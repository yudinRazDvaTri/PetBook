package authentication

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	gorillaContext "github.com/gorilla/context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"time"
)

// Structure which represents jwt claims.
type Claims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

// Simply pair of access and refresh tokens and their expiration time.
type Tokens struct {
	AccessTokenValue      string
	AccessExpirationTime  time.Time
	RefreshTokenValue     string
	RefreshExpirationTime time.Time
}

// This middleware restricts users, who haven't logged in (haven't entered email and password), from endpoints.
// Firstly, it checks whether user has logged in with the help of 3-rd party services or not.
// If not, then access token is checked. If there is no access token, then refresh token is checked.
// If none of the options fit, then the user is not authenticated and needs to be redirected to the login page.
func AuthenticateMiddleware(storeRefreshToken *models.RefreshTokenStore) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Checking whether client has a cookie, which is called "oauth".
			oauthContent, err := r.Cookie("oauth")
			if err == nil {
				// Decoding and storing cookie value into the 'value' variable.
				value := make(map[string]interface{})
				if err = SCookie.Decode("oauth", oauthContent.Value, &value); err == nil {
					oauthToken := value["accessToken"].(*oauth2.Token)

					// Getting actual token pair from Google.
					// TokenSource is anything that can return a token and
					// it holds new pair of access and refresh tokens.
					tokenSource := GoogleOauthConfig.TokenSource(context.Background(), oauthToken)

					// Selecting pair of tokens from TokenSource.
					newToken, err := tokenSource.Token()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logger.Error(err, "Error occurred while trying to get Token from TokenSource.\n")
						return
					}

					// If life time of 3-rd party access token expires, we will use new pair of tokens.
					if newToken.AccessToken != oauthToken.AccessToken {
						client := oauth2.NewClient(context.Background(), tokenSource)

						// Requesting content about user from Google resource server with the help of access token.
						response, err := client.Get(OauthGoogleUrlAPI + newToken.AccessToken)
						if err != nil {
							logger.Error()
							http.Redirect(w, r, "/login", http.StatusSeeOther)
							return
						}

						contents, err := ioutil.ReadAll(response.Body)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							logger.Error(err, "Error occurred while trying to read user info bytes.\n")
							return
						}

						var googleUserInfo GoogleUserInfo

						if err := json.Unmarshal(contents, &googleUserInfo); err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							logger.Error(err, "Error occurred while trying to unmarshal user info.\n")
							return
						}

						gob.Register(newToken)
						value := map[string]interface{}{
							"accessToken": newToken,
							"userId":      value["userId"].(int),
						}

						// Setting new tokens into a cookie.
						if encoded, err := SCookie.Encode("oauth", value); err == nil {
							cookie := &http.Cookie{
								Name:  "oauth",
								Value: encoded,

								// Expiration time of cookie which stores oauth information was set twice
								// as much as google oauth token expiration time.
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
					}

					gorillaContext.Set(r, "id", value["userId"].(int))
					h.ServeHTTP(w, r)
					return

				} else {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
			}

			// Checking whether client has a cookie, which is called "accessToken".
			accessToken, err := r.Cookie("accessToken")
			if err != nil {
				// If not checking for "refreshToken".
				refreshToken, err := r.Cookie("refreshToken")
				if err != nil {
					// Redirect to login page if none of the options fit.
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}

				refreshTokenString := refreshToken.Value
				claims := &Claims{}

				// Validating refresh token.
				token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return RsaKeys.VerifyKey, nil
				})

				if err == nil {
					// Checking whether the same client tries to get new pair of tokens.
					userAgent := GetUserAgent(r)
					err := storeRefreshToken.RefreshTokenExists(claims.Id, refreshTokenString, userAgent)
					if err != nil {
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

					// If refresh token is valid, then new pair of tokens is generated.
					if token.Valid {
						tokens, err := GenerateTokenPair(claims.Id)
						if err != nil {
							logger.Error(err)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						http.SetCookie(w, &http.Cookie{
							Name:    "accessToken",
							Value:   tokens.AccessTokenValue,
							Expires: tokens.AccessExpirationTime,
							Path:    "/",
						})

						// Updating refresh token in database.
						err = storeRefreshToken.UpdateRefreshToken(claims.Id, tokens.RefreshTokenValue, tokens.RefreshExpirationTime, userAgent)
						if err != nil {
							logger.Error(err)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						http.SetCookie(w, &http.Cookie{
							Name:    "refreshToken",
							Value:   tokens.RefreshTokenValue,
							Expires: tokens.RefreshExpirationTime,
							Path:    "/",
						})

						gorillaContext.Set(r, "id", claims.Id)

					} else {
						http.Redirect(w, r, "/login", http.StatusSeeOther)
						return
					}
				} else {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}

			} else {
				accessTokenString := accessToken.Value
				claims := &Claims{}

				// Validating access token
				token, err := jwt.ParseWithClaims(accessTokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return RsaKeys.VerifyKey, nil
				})

				if err == nil {
					if token.Valid {
						gorillaContext.Set(r, "id", claims.Id)
					} else {
						http.Redirect(w, r, "/login", http.StatusSeeOther)
						return
					}
				} else {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
	return
}

// This middleware allows user to request endpoints if he chose a role.
// In other way, user will be redirected to role page.
// This middleware is used with AuthenticateMiddleware.
func AuthorizeMiddleware(storeUser *models.UserStore) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userId := gorillaContext.Get(r, "id").(int)
			role, _ := storeUser.GetUserRole(userId)
			if role == "" {
				http.Redirect(w, r, "/role", http.StatusSeeOther)
				return
			}
			gorillaContext.Set(r, "role", role)
			h.ServeHTTP(w, r)
		})
	}
	return
}

// PetOrVetMiddleware restricts users, who haven't registered their entities (pet or vet), from endpoints.
// This middleware is used with AuthenticateMiddleware and AuthorizeMiddleware.
// PetOrVetMiddleware is the maximum degree of authorization in this application.
func PetOrVetMiddleware(storeUser *models.UserStore) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := gorillaContext.Get(r, "id").(int)
			role := gorillaContext.Get(r, "role").(string)

			switch role {
			case "pet":
				_, err := storeUser.GetPet(userId)
				if err != nil {
					http.Redirect(w, r, "/petcabinet", http.StatusSeeOther)
					return
				}
			case "vet":
				_, err := storeUser.GetVet(userId)
				if err != nil {
					http.Redirect(w, r, "/vetcabinet", http.StatusSeeOther)
					return
				}
			default:
				http.Error(w, "Unknown user role.", http.StatusBadRequest)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
	return
}

// Generating new pair of access and refresh token when user is logging in or when tokens expire.
func GenerateTokenPair(userID int) (Tokens, error) {
	var tokens Tokens

	tokens.AccessExpirationTime = time.Now().Add(5 * time.Minute)
	accessClaims := &Claims{
		Id: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokens.AccessExpirationTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(RsaKeys.SignKey)

	if err != nil {
		return tokens, fmt.Errorf("Error occurred while trying to sign access token: %v.\n", err)
	}
	tokens.AccessTokenValue = accessTokenString

	tokens.RefreshExpirationTime = time.Now().Add(60 * 24 * time.Hour)
	refreshClaims := &Claims{
		Id: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokens.RefreshExpirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(RsaKeys.SignKey)

	if err != nil {
		return tokens, fmt.Errorf("Error occurred while trying to sign refresh token: %v.\n", err)
	}

	tokens.RefreshTokenValue = refreshTokenString

	return tokens, nil
}

func GetUserAgent(r *http.Request) string {
	ua := r.Header.Get("User-Agent")
	return ua
}
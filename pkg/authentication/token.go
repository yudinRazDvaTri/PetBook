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
	_ "golang.org/x/oauth2"
	"io/ioutil"
	_ "io/ioutil"

	"net/http"
	"time"
)

type Claims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

type Tokens struct {
	AccessTokenValue      string
	AccessExpirationTime  time.Time
	RefreshTokenValue     string
	RefreshExpirationTime time.Time
}

func AuthMiddleware(storeRefreshToken *models.RefreshTokenStore) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Checking whether user logs in with the help of third-party service
			oauthContent, err := r.Cookie("oauth")
			if err == nil {
				value := make(map[string]interface{})
				if err = SCookie.Decode("oauth", oauthContent.Value, &value); err == nil {
					oauthToken := value["accessToken"].(*oauth2.Token)
					tokenSource := GoogleOauthConfig.TokenSource(context.Background(), oauthToken)

					newToken, err := tokenSource.Token()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logger.Error( err, "Error occurred while trying to get Token from TokenSource.\n")
						return
					}

					if newToken.AccessToken != oauthToken.AccessToken {
						client := oauth2.NewClient(context.Background(), tokenSource)
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

						if encoded, err := SCookie.Encode("oauth", value); err == nil {
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
							logger.Error(err.Error(), "; Error occurred while trying to encode cookie.\n", )
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

			accessToken, err := r.Cookie("accessToken")
			if err != nil {
				refreshToken, err := r.Cookie("refreshToken")
				if err != nil {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}

				refreshTokenString := refreshToken.Value
				claims := &Claims{}

				// Validate refresh token
				token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return RsaKeys.VerifyKey, nil
				})

				if err == nil {
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

				// Validate access token
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

func PetMiddleware(storeUser *models.UserStore) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := gorillaContext.Get(r, "id").(int)

			_, err := storeUser.GetPet(userId)
			if err != nil {
				http.Redirect(w, r, "/petcabinet", http.StatusSeeOther)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
	return
}

// Generating new pair of access and refresh token when user is logging or when tokens expire
func GenerateTokenPair(userID int) (Tokens, error) {
	var tokens Tokens

	tokens.AccessExpirationTime = time.Now().Add(1 * time.Minute)
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
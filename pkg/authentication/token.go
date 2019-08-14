package authentication

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/dpgolang/PetBook/pkg/models"
	"github.com/dpgolang/PetBook/pkg/utilerr"
	"github.com/gorilla/context"
	_ "github.com/lib/pq"
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

func ValidateTokenMiddleware(storeRefreshToken *models.RefreshTokenStore) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := r.Cookie("accessToken")
			if err != nil {
				refreshToken, err := r.Cookie("refreshToken")
				if err != nil {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}

				refreshTokenString := refreshToken.Value
				claims := &Claims{}

				//validate refresh token
				token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return Keys.VerifyKey, nil
				})

				if err == nil {
					err := storeRefreshToken.RefreshTokenExists(claims.Id, refreshTokenString)
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
						})

						err = storeRefreshToken.UpdateRefreshToken(claims.Id, tokens.RefreshTokenValue, tokens.RefreshExpirationTime)
						if err != nil {
							logger.Error(err)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						http.SetCookie(w, &http.Cookie{
							Name:    "refreshToken",
							Value:   tokens.RefreshTokenValue,
							Expires: tokens.RefreshExpirationTime,
						})
						context.Set(r, "id", claims.Id)
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

				//validate access token
				token, err := jwt.ParseWithClaims(accessTokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return Keys.VerifyKey, nil
				})

				if err == nil {
					if token.Valid {
						context.Set(r, "id", claims.Id)
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
	accessTokenString, err := accessToken.SignedString(Keys.SignKey)

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
	refreshTokenString, err := refreshToken.SignedString(Keys.SignKey)

	if err != nil {
		return tokens, fmt.Errorf("Error occurred while trying to sign refresh token: %v.\n", err)
	}

	tokens.RefreshTokenValue = refreshTokenString

	return tokens, nil
}
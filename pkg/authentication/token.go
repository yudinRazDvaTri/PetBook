package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	_ "github.com/lib/pq"
	"net/http"
)

type Claims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := c.Value
	claims := &Claims{}

	//r.Header.Set("Authorization", "Bearer " + tokenString)

	/*token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
	func(token *jwt.Token) (interface{}, error) {
		return VerifyKey, nil
	})
	*/

	//validate token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Keys.VerifyKey, nil
	})

	if err == nil {
		if token.Valid {
			context.Set(r, "id", claims.Id)
			next(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func Content(h http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString := c.Value
		claims := &Claims{}

		//r.Header.Set("Authorization", "Bearer " + tokenString)

		/*token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return VerifyKey, nil
		})
		*/

		//validate token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return Keys.VerifyKey, nil
		})

		if err == nil {
			if token.Valid {
				context.Set(r, "id", claims.Id)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		h.ServeHTTP(w, r)
	})
}

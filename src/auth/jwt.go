package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// JWT json web token
func (a *Auth) JWT(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			token, err := jwt.ParseWithClaims(auth, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(a.Key), nil
			})
			if err == nil {
				if tk, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
					r.Header.Set("uid", tk.Audience)
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

// GenJWT generate jwt
func (a *Auth) GenJWT(aud string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: 60 * 60 * 24 * 15,
		Issuer:    "fuRan",
		Audience:  aud,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.Key)
}

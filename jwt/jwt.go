package jwt

import (
	"api/internal/util"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(handlerFnWithJWT http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-hunter999-token")
		_, err := jwtValidator(tokenString)
		if err != nil {
			util.WriteJSONResponse(w, http.StatusForbidden, map[string]string{"error": "invalid token"})
			return
		}
		handlerFnWithJWT(w, r)
	}
}

func jwtValidator(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

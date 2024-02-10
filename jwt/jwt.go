package midd

import (
	"api/internal/util"
	"api/store"
	"api/types"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(handlerFnWithJWT http.HandlerFunc, s store.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-hunter999-token")
		token, err := jwtValidator(tokenString)
		if err != nil || !token.Valid {
			util.WriteJSONResponse(w, http.StatusForbidden, util.NewError(types.ErrAccessDenied))
			return
		}
		userId, err := util.IDGetter(r)
		if err != nil {
			util.WriteJSONResponse(w, http.StatusForbidden, util.NewError(types.ErrAccessDenied))
			return
		}
		account, err := s.GetAccountById(userId)
		if err != nil {
			util.WriteJSONResponse(w, http.StatusForbidden, util.NewError(types.ErrAccessDenied))
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if account.Number != claims["accountNumber"] {
			util.WriteJSONResponse(w, http.StatusForbidden, util.NewError(types.ErrAccessDenied))
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

func NewJWTToken(account *types.Account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

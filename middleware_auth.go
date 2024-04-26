package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loyalsfc/fledge-backend/internal/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, string)

func (apiCfg *apiCfg) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			errResponse(403, w, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		token, err := jwt.Parse(apiKey, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			errResponse(403, w, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		username := claims["username"]

		handler(w, r, fmt.Sprintf("%v", username))
	}
}

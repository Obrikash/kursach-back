package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/obrikash/swimming_pool/internal/data"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		tokenString := headerParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
			}

			return []byte(app.config.jwt.secret), nil
		})
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			now := time.Now()

			exp := claims["exp"].(float64)
			expTime := time.Unix(int64(exp), 0)

			nbf := claims["nbf"].(float64)
			nbfTime := time.Unix(int64(nbf), 0)
			if now.After(expTime) || now.Before(nbfTime) {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
		} else {
			app.serverErrorResponse(w, r, err)
			return
		}
		if claims["iss"].(string) != "github.com/obrikash/swimming_pool" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		user, err := app.models.Users.Get(userID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)

	})
}

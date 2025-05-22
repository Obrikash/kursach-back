package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.HandlerFunc(http.MethodGet, "/v1/pools", app.listPoolsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/pool", app.mostProfitPoolHandler)

	router.HandlerFunc(http.MethodGet, "/v1/users/trainers", app.listTrainersHandler)
	router.HandlerFunc(http.MethodGet, "/v1/pools/trainers", app.listTrainersForPoolsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/trainers/profit", app.profitOfTrainers)

	router.HandlerFunc(http.MethodGet, "/v1/groups", app.listGroupsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/groups", app.addGroupToPoolHandler)

	router.HandlerFunc(http.MethodGet, "/v1/subscriptions", app.listSubscriptionsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	return app.recoverPanic(app.authenticate(router))
}

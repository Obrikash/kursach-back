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

	router.HandlerFunc(http.MethodGet, "/v1/users/trainers", app.requireAuthenticatedUser(app.listTrainersHandler))
	router.HandlerFunc(http.MethodGet, "/v1/pools/trainers", app.requireAuthenticatedUser(app.listTrainersForPoolsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/pools/trainers", app.requireAuthenticatedUser(app.attachTrainerToPoolHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/trainers/profit", app.requireAuthenticatedUser(app.profitOfTrainers))

	router.HandlerFunc(http.MethodGet, "/v1/groups", app.requireAuthenticatedUser(app.listGroupsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/groups", app.requireAuthenticatedUser(app.addGroupToPoolHandler))

	router.HandlerFunc(http.MethodGet, "/v1/subscriptions", app.requireAuthenticatedUser(app.listSubscriptionsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	return app.recoverPanic(app.enableCORS((app.authenticate(router))))
}

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/pools", app.listPoolsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/trainers", app.listTrainersHandler)
	router.HandlerFunc(http.MethodGet, "/v1/groups", app.listGroupsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/subscriptions", app.listSubscriptionsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/pools/trainers", app.listTrainersForPoolsHandler)
	return router
}

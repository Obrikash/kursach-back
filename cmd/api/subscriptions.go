package main

import "net/http"

func (app *application) listSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := app.models.Subscriptions.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"subscriptions": subscriptions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listUsersSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	subscriptions, err := app.models.Subscriptions.UserSubscriptions(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user_subscriptions": subscriptions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

package main

import "net/http"

func (app *application) listGroupsHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := app.models.Groups.GetGroups()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"groups": groups}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
